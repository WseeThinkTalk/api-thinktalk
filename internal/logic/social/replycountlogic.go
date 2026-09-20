package social

import (
	"context"

	reply "api-thinktalk/client/reply/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCountLogic {
	return &ReplyCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyCountLogic) ReplyCount(req *types.ReplyCountRequest) (*types.ReplyCountResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyCount(l.ctx, &reply.ReplyCountRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[ReplyCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCountResponse{
		ReplyNum:     resp.ReplyNum,
		ReplyRootNum: resp.ReplyRootNum,
	}, nil
}
