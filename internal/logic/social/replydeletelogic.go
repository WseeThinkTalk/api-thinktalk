package social

import (
	"context"

	reply "api-thinktalk/client/reply/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyDeleteLogic {
	return &ReplyDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyDeleteLogic) DeleteReply(userId int64, req *types.ReplyDeleteRequest) (*types.ReplyDeleteResponse, error) {
	_, err := l.svcCtx.ReplyRPC.DeleteReply(l.ctx, &reply.DeleteReplyRequest{
		ReplyId: req.ReplyId,
		UserId:  userId,
	})
	if err != nil {
		l.Errorf("[DeleteReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyDeleteResponse{}, nil
}
