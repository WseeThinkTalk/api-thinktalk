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
		TargetId: req.ObjId,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}

	if rpcResp != nil && rpcResp.Data != nil && rpcResp.Data.UserThumbups != nil {
		if thumbup, ok := rpcResp.Data.UserThumbups[req.ObjId]; ok && thumbup != nil {
			resp.HasLiked = thumbup.LikeType == 1
		}
	}
	return resp, nil
}
