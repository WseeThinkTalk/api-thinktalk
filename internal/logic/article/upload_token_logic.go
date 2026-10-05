package article

import (
	"context"
	"fmt"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// UploadTokenLogic 获取文件上传凭证
type UploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUploadTokenLogic 初始化逻辑对象
func NewUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTokenLogic {
	return &UploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadToken 生成预签名上传 URL
func (l *UploadTokenLogic) UploadToken(req *types.UploadTokenRequest) (resp *types.UploadTokenResponse, err error) {
	resp = new(types.UploadTokenResponse)

	// 校验文件格式
	ext, ok := allowedMimeTypes[req.FileType]
	if !ok {
		return nil, fmt.Errorf("不支持的文件类型，仅允许上传 jpg/png/webp/gif 图片")
	}

	// 业务场景目录
	scene := req.Scene
	if scene != "cover" && scene != "avatar" && scene != "article" {
		scene = "common"
	}

	// 生成对象路径与过期时间（15分钟）
	objectKey := fmt.Sprintf("%s/%s%s", scene, uuid.New().String(), ext)
	expires := 15 * time.Minute

	// 生成预签名 URL
	presignedURL, err := l.svcCtx.MinIO.PresignedPutObject(
		l.ctx,
		l.svcCtx.Config.MinIO.BucketName,
		objectKey,
		expires,
	)
	if err != nil {
		l.Errorf("生成预签名 URL 失败: %v", err)
		return nil, err
	}

	// 组装返回结果
	resp.UploadUrl = presignedURL.String()
	resp.FileKey = objectKey
	resp.ExpireSec = int64(expires.Seconds())
	resp.ViewUrl = fmt.Sprintf("/static/%s/%s", l.svcCtx.Config.MinIO.BucketName, objectKey)

	return resp, nil
}
