package video

import (
	"context"
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

// UploadInit 初始化分片上传并下发直传预签名URL
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

	ext := allowedVideoMimes[req.FileType]
	objectKey := fmt.Sprintf("video/raw/%s/%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)
	bucket := l.svcCtx.Config.MinIO.BucketName

	// 1. 调用 MinIO Core 初始化 S3 分片上传
	core := minio.Core{Client: l.svcCtx.MinIO}
	uploadID, err := core.NewMultipartUpload(l.ctx, bucket, objectKey, minio.PutObjectOptions{
		ContentType: req.FileType,
	})
	if err != nil {
		l.Errorf("初始化分片上传失败: %v", err)
		return nil, fmt.Errorf("初始化分片上传失败")
	}

	// 2. 批量签发每个分片的 Presigned PUT URL
	expires := 2 * time.Hour
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

	videoId := snowflake.GenerateID()

	// 3. 将上传元数据写入 Redis 缓存状态表（TTL: 2小时）
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

	return resp, nil
}
