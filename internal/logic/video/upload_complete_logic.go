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

// UploadComplete 通知 MinIO 完成合并并投递流媒体处理消息
func (l *UploadCompleteLogic) UploadComplete(req *types.VideoUploadCompleteRequest) (resp *types.VideoUploadCompleteResponse, err error) {
	resp = new(types.VideoUploadCompleteResponse)

	bucket := l.svcCtx.Config.MinIO.BucketName
	core := minio.Core{Client: l.svcCtx.MinIO}

	// 1. 组装 CompletePart 切片
	parts := make([]minio.CompletePart, len(req.Parts))
	for i, p := range req.Parts {
		parts[i] = minio.CompletePart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		}
	}

	// 2. 通知 MinIO 完成分片合并
	_, err = core.CompleteMultipartUpload(l.ctx, bucket, req.ObjectKey, req.UploadId, parts, minio.PutObjectOptions{})
	if err != nil {
		l.Errorf("合并分片失败: %v", err)
		return nil, fmt.Errorf("合并分片失败: %w", err)
	}

	viewUrl := fmt.Sprintf("/static/%s/%s", bucket, req.ObjectKey)

	// 3. 更新 Redis 状态为 processing
	if l.svcCtx.BizRedis != nil {
		redisKey := fmt.Sprintf("biz#video#status:%d", req.VideoId)
		statusVal := fmt.Sprintf(`{"status":"processing","videoUrl":"%s","objectKey":"%s"}`, viewUrl, req.ObjectKey)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, redisKey, statusVal, 86400*7)

		// 建立哈希到已完成视频的全局映射索引，后续全网相同文件直接命中秒传（TTL: 30天）
		if req.FileHash != "" {
			hashKey := fmt.Sprintf("biz#video#hash:%s", req.FileHash)
			hashVal := fmt.Sprintf(`{"status":"ready","videoId":%d,"objectKey":"%s","videoUrl":"%s"}`,
				req.VideoId, req.ObjectKey, viewUrl)
			_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, hashKey, hashVal, 86400*30)
		}
	}

	// 4. 投递异步流媒体处理消息至 Kafka
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
