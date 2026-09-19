package logic

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"api-thinktalk/article/internal/code"
	"api-thinktalk/article/internal/svc"
	"api-thinktalk/article/internal/types"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxFileSize   = 10 << 20
	maxImagePixels = 89478485 // ~8K resolution (8192*8192 + safety margin)
)

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	if err := req.ParseMultipartForm(maxFileSize); err != nil {
		return nil, code.ParseFormErr
	}

	file, header, err := req.FormFile("cover")
	if err != nil {
		logx.Errorf("get form file failed, err: %v", err)
		return nil, err
	}
	defer file.Close()

	// Read first 512 bytes for content-type detection
	buf := make([]byte, 512)
	n, _ := io.ReadFull(file, buf)
	if n == 0 {
		return nil, fmt.Errorf("empty file")
	}
	buf = buf[:n]

	detectedType := http.DetectContentType(buf)
	ext, ok := allowedMimeTypes[detectedType]
	if !ok {
		logx.Errorf("rejected file type: %s (detected), client claimed: %s", detectedType, header.Header.Get("Content-Type"))
		return nil, fmt.Errorf("不支持的文件类型，仅允许 jpg/png/webp/gif")
	}

	// Validate image dimensions to prevent decompression bombs
	imgConfig, _, err := image.DecodeConfig(bytes.NewReader(buf))
	if err == nil {
		pixels := imgConfig.Width * imgConfig.Height
		if pixels > maxImagePixels {
			return nil, fmt.Errorf("图片分辨率过大，请上传小于 8K 的图片")
		}
	}

	// Build a reader that includes the bytes we already consumed
	_ = req.MultipartForm.RemoveAll()
	req.MultipartForm = nil

	// Re-read the file since we consumed the first 512 bytes
	file.Seek(0, io.SeekStart)

	objectKey := fmt.Sprintf("cover/%s%s", uuid.New().String(), ext)
	_, err = l.svcCtx.MinIO.PutObject(
		l.ctx,
		l.svcCtx.Config.MinIO.BucketName,
		objectKey,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: detectedType,
		},
	)
	if err != nil {
		logx.Errorf("put object to minio failed, err: %v", err)
		return nil, code.PutBucketErr
	}
	return &types.UploadCoverResponse{
		CoverUrl: genFileURL(
			l.svcCtx.Config.MinIO.Endpoint,
			l.svcCtx.Config.MinIO.BucketName,
			objectKey)}, nil
}

func genFileURL(endpoint, bucketName, objectKey string) string {
	// Only expose the bucket/object path, not the raw MinIO endpoint
	return fmt.Sprintf("/static/%s/%s", bucketName, objectKey)
}
