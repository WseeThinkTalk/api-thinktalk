package video

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/snowflake"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

var allowedVideoMimes = map[string]string{
	"video/mp4":        ".mp4",
	"video/quicktime":  ".mov",
	"video/webm":       ".webm",
	"video/x-matroska": ".mkv",
}

// IsValidVideoMimeType 检查是否是受支持的视频MIME格式
func IsValidVideoMimeType(mime string) bool {
	_, ok := allowedVideoMimes[mime]
	return ok
}

// CalculatePartCount 计算分片数量
func CalculatePartCount(fileSize, partSize int64) int {
	if partSize <= 0 {
		partSize = 10 * 1024 * 1024 // 默认 10MB
	}
	parts := fileSize / partSize
	if fileSize%partSize != 0 {
		parts++
	}
	if parts == 0 {
		parts = 1
	}
	return int(parts)
}

// UploadInitLogic 视频分片初始化逻辑
type UploadInitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUploadInitLogic 实例化分片初始化逻辑
func NewUploadInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadInitLogic {
	return &UploadInitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadInit 初始化分片上传、支持秒传判定与断点续传探测并下发直传预签名URL
func (l *UploadInitLogic) UploadInit(req *types.VideoUploadInitRequest) (resp *types.VideoUploadInitResponse, err error) {
	resp = new(types.VideoUploadInitResponse)

	if !IsValidVideoMimeType(req.FileType) {
		return nil, fmt.Errorf("不支持的视频格式，仅支持 mp4/mov/webm")
	}

	partSize := req.PartSize
	if partSize <= 0 {
		partSize = 10 * 1024 * 1024 // 默认 10MB
	}
	partCount := CalculatePartCount(req.FileSize, partSize)
	bucket := l.svcCtx.Config.MinIO.BucketName
	expires := 2 * time.Hour

	var (
		videoId       int64
		uploadID      string
		objectKey     string
		isResuming    bool
		uploadedParts []int
	)

	core := minio.Core{Client: l.svcCtx.MinIO}

	// 1. 秒传与断点续传探测（若客户端提供了文件全量 SHA-256 哈希）
	if req.FileHash != "" && l.svcCtx.BizRedis != nil {
		hashKey := fmt.Sprintf("biz#video#hash:%s", req.FileHash)
		hashVal, err := l.svcCtx.BizRedis.GetCtx(l.ctx, hashKey)
		if err == nil && hashVal != "" {
			var cached struct {
				Status    string `json:"status"`
				VideoID   int64  `json:"videoId"`
				UploadID  string `json:"uploadId"`
				ObjectKey string `json:"objectKey"`
			}
			if json.Unmarshal([]byte(hashVal), &cached) == nil {
				// 1.1 命中秒传：全网已存在相同视频且已就绪，零流量瞬时完成！
				if cached.Status == "ready" {
					resp.VideoId = cached.VideoID
					resp.ObjectKey = cached.ObjectKey
					resp.IsQuickDone = true
					resp.PartSize = partSize
					resp.PartCount = partCount
					resp.ExpireSec = int64(expires.Seconds())
					l.Infof("[UploadInit] Instant upload matched for hash %s, videoId: %d", req.FileHash, cached.VideoID)
					return resp, nil
				}

				// 1.2 断点续传探测：正在上传中，探测已成功上传的分片编号列表
				if cached.Status == "uploading" && cached.UploadID != "" && cached.ObjectKey != "" {
					partsResult, err := core.ListObjectParts(l.ctx, bucket, cached.ObjectKey, cached.UploadID, 0, 1000)
					if err == nil {
						uploadedParts = make([]int, 0, len(partsResult.ObjectParts))
						for _, p := range partsResult.ObjectParts {
							uploadedParts = append(uploadedParts, p.PartNumber)
						}
						videoId = cached.VideoID
						uploadID = cached.UploadID
						objectKey = cached.ObjectKey
						isResuming = true
						l.Infof("[UploadInit] Breakpoint resume detected for hash %s, videoId: %d, parts: %v",
							req.FileHash, videoId, uploadedParts)
					}
				}
			}
		}
	}

	// 2. 若未复用断点会话，则创建全新的 MinIO Multipart Upload 与 Snowflake 视频ID
	if !isResuming {
		ext := allowedVideoMimes[req.FileType]
		objectKey = fmt.Sprintf("video/raw/%s/%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)
		newUploadID, err := core.NewMultipartUpload(l.ctx, bucket, objectKey, minio.PutObjectOptions{
			ContentType: req.FileType,
		})
		if err != nil {
			l.Errorf("初始化分片上传失败: %v", err)
			return nil, fmt.Errorf("初始化分片上传失败")
		}
		uploadID = newUploadID
		videoId = snowflake.GenerateID()

		// 若携带了哈希，记录 uploading 状态供中断重试续传
		if req.FileHash != "" && l.svcCtx.BizRedis != nil {
			hashKey := fmt.Sprintf("biz#video#hash:%s", req.FileHash)
			metaJSON, _ := json.Marshal(map[string]interface{}{
				"status":    "uploading",
				"videoId":   videoId,
				"uploadId":  uploadID,
				"objectKey": objectKey,
			})
			_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, hashKey, string(metaJSON), int(expires.Seconds()))
		}
	}

	// 3. 批量签发每个分片的 Presigned PUT 直传 URL
	partUrls := make([]string, partCount)
	for i := 1; i <= partCount; i++ {
		queryParams := make(url.Values)
		queryParams.Set("uploadId", uploadID)
		queryParams.Set("partNumber", strconv.Itoa(i))

		partURL, err := l.svcCtx.MinIO.Presign(l.ctx, http.MethodPut, bucket, objectKey, expires, queryParams)
		if err != nil {
			l.Errorf("签发分片 URL 失败 [part %d]: %v", i, err)
			return nil, fmt.Errorf("签发分片凭证失败")
		}
		partUrls[i-1] = partURL.String()
	}

	// 4. 将上传元数据写入 Redis 缓存状态表（TTL: 2小时）
	if l.svcCtx.BizRedis != nil {
		redisKey := fmt.Sprintf("biz#video#status:%d", videoId)
		videoData := fmt.Sprintf(`{"status":"uploading","objectKey":"%s","uploadId":"%s"}`, objectKey, uploadID)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, redisKey, videoData, int(expires.Seconds()))
	}

	resp.VideoId = videoId
	resp.UploadId = uploadID
	resp.ObjectKey = objectKey
	resp.PartUrls = partUrls
	resp.PartSize = partSize
	resp.PartCount = partCount
	resp.ExpireSec = int64(expires.Seconds())
	resp.IsQuickDone = false
	resp.UploadedParts = uploadedParts

	return resp, nil
}
