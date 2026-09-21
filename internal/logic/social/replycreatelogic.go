package social

import (
	"context"

	reply "api-thinktalk/client/reply/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCreateLogic {
	return &ReplyCreateLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyCreateLogic) CreateReply(userId int64, req *types.ReplyCreateRequest) (resp *types.ReplyCreateResponse, err error) {
	resp = new(types.ReplyCreateResponse)

	rpcResp, err := l.svcCtx.ReplyRPC.CreateReply(l.ctx, &reply.CreateReplyRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   userId,
		BeReplyUserId: req.BeReplyUserId,
		ParentId:      req.ParentId,
		Content:       req.Content,
	})
	if err != nil {
		l.Errorf("[CreateReply] rpc err: %v", err)
		return nil, err
	}
	resp.ReplyId = rpcResp.ReplyId
	return resp, nil
}
