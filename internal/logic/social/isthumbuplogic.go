package social

import (
	"context"
	"encoding/json"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	like "api-thinktalk/client/like/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsThumbupLogic) IsThumbup(req *types.IsThumbupRequest) (resp *types.IsThumbupResponse, err error) {
	resp = new(types.IsThumbupResponse)

	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.LikeRPC.IsThumbup(l.ctx, &like.IsThumbupRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}

	if thumbup, ok := rpcResp.UserThumbups[req.TargetId]; ok {
		resp.LikeType = thumbup.LikeType
		resp.ThumbupTime = thumbup.ThumbupTime
	}
	return resp, nil
}
