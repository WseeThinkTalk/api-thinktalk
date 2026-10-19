package video

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

// UploadCompleteLogic 视频分片合并与事件触发逻辑
type UploadCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUploadCompleteLogic 实例化分片合并逻辑
func NewUploadCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCompleteLogic {
	return &UploadCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadComplete 完成分片合并并投递转码处理消息
func (l *UploadCompleteLogic) UploadComplete(req *types.VideoUploadCompleteRequest) (resp *types.VideoUploadCompleteResponse, err error) {
	resp = new(types.VideoUploadCompleteResponse)

	if len(req.Parts) == 0 {
		return nil, fmt.Errorf("分片列表不能为空")
	}

	// 状态检查与权限校验
	if l.svcCtx.BizRedis != nil {
		redisKey := fmt.Sprintf("biz#video#status:%d", req.VideoId)
		val, err := l.svcCtx.BizRedis.GetCtx(l.ctx, redisKey)
		if err == nil && val != "" {
			var meta struct {
				Status   string `json:"status"`
				UserId   int64  `json:"userId"`
				VideoUrl string `json:"videoUrl"`
			}
			if json.Unmarshal([]byte(val), &meta) == nil {
				// 校验归属，防止越权
				if meta.UserId > 0 {
					if uidVal := l.ctx.Value("userId"); uidVal != nil {
						if n, ok := uidVal.(json.Number); ok {
							if currUid, _ := n.Int64(); currUid > 0 && currUid != meta.UserId {
								return nil, fmt.Errorf("无权操作此视频")
							}
						}
					}
				}
				// 已处理则直接返回，避免重复调 MinIO 和发 MQ
				if meta.Status == "processing" || meta.Status == "ready" {
					resp.VideoId = req.VideoId
					resp.Status = meta.Status
					resp.VideoUrl = meta.VideoUrl
					return resp, nil
				}
			}
		}
	}

	bucket := l.svcCtx.Config.MinIO.BucketName
	core := minio.Core{Client: l.svcCtx.MinIO}

	// 组装分片列表
	parts := make([]minio.CompletePart, len(req.Parts))
	for i, p := range req.Parts {
		parts[i] = minio.CompletePart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		}
	}

	// 调用 MinIO 合并分片
	_, err = core.CompleteMultipartUpload(l.ctx, bucket, req.ObjectKey, req.UploadId, parts, minio.PutObjectOptions{})
	if err != nil {
		l.Errorf("合并分片失败: %v", err)
		return nil, fmt.Errorf("合并分片失败: %w", err)
	}

	viewUrl := fmt.Sprintf("/static/%s/%s", bucket, req.ObjectKey)

	// 更新处理状态
	if l.svcCtx.BizRedis != nil {
		redisKey := fmt.Sprintf("biz#video#status:%d", req.VideoId)
		statusVal := fmt.Sprintf(`{"status":"processing","videoUrl":"%s","objectKey":"%s"}`, viewUrl, req.ObjectKey)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, redisKey, statusVal, 86400*7)

		// 记录文件哈希映射，供后续秒传复用（TTL: 30天）
		if req.FileHash != "" {
			hashKey := fmt.Sprintf("biz#video#hash:%s", req.FileHash)
			hashVal := fmt.Sprintf(`{"status":"ready","videoId":%d,"objectKey":"%s","videoUrl":"%s"}`,
				req.VideoId, req.ObjectKey, viewUrl)
			_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, hashKey, hashVal, 86400*30)
		}
	}

	// 投递转码处理消息至 Kafka
	if l.svcCtx.NotificationPusher != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"videoId":   req.VideoId,
			"objectKey": req.ObjectKey,
			"bucket":    bucket,
			"videoUrl":  viewUrl,
			"timestamp": time.Now().Unix(),
		})
		_ = l.svcCtx.NotificationPusher.Push(l.ctx, string(eventData))
	}

	resp.VideoId = req.VideoId
	resp.Status = "processing"
	resp.VideoUrl = viewUrl

	return resp, nil
}
