package video

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// VideoStatusLogic 视频处理就绪状态轮询逻辑
type VideoStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewVideoStatusLogic 实例化状态轮询逻辑
func NewVideoStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoStatusLogic {
	return &VideoStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VideoStatus 查询视频处理状态及抽帧封面与元数据
func (l *VideoStatusLogic) VideoStatus(req *types.VideoStatusRequest) (resp *types.VideoStatusResponse, err error) {
	resp = new(types.VideoStatusResponse)
	resp.VideoId = req.VideoId
	resp.Status = "processing"

	if l.svcCtx.BizRedis == nil {
		return resp, nil
	}

	redisKey := fmt.Sprintf("biz#video#status:%d", req.VideoId)
	val, err := l.svcCtx.BizRedis.GetCtx(l.ctx, redisKey)
	if err != nil || val == "" {
		resp.Status = "not_found"
		return resp, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return resp, nil
	}

	if s, ok := data["status"].(string); ok {
		resp.Status = s
	}
	if vu, ok := data["videoUrl"].(string); ok {
		resp.VideoUrl = vu
	}
	if cu, ok := data["coverUrl"].(string); ok {
		resp.CoverUrl = cu
	}
	if dur, ok := data["duration"].(float64); ok {
		resp.Duration = int64(dur)
	}
	if w, ok := data["width"].(float64); ok {
		resp.Width = int(w)
	}
	if h, ok := data["height"].(float64); ok {
		resp.Height = int(h)
	}
	if em, ok := data["errMsg"].(string); ok {
		resp.ErrMsg = em
	}

	return resp, nil
}
