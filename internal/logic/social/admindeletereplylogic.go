package social

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	reply "api-thinktalk/client/reply/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminDeleteReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteReplyLogic {
	return &AdminDeleteReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminDeleteReplyLogic) AdminDeleteReply(req *types.AdminDeleteReplyRequest) (*types.AdminDeleteReplyResponse, error) {
	_, err := l.svcCtx.ReplyRPC.DeleteReply(l.ctx, &reply.DeleteReplyRequest{
		ReplyId: req.ReplyId,
		IsAdmin: true,
	})
	if err != nil {
		l.Errorf("[AdminDeleteReply] rpc err: %v replyId: %d", err, req.ReplyId)
		return nil, err
	}
	return &types.AdminDeleteReplyResponse{}, nil
}
