package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkReadLogic) MarkRead(userId int64, req *types.ChatMarkReadRequest) (resp *types.ChatMarkReadResponse, err error) {
	resp = new(types.ChatMarkReadResponse)

	_, err = l.svcCtx.Chat.MarkRead(l.ctx, &pb.ChatMarkReadRequest{
		UserId:         userId,
		ConversationId: req.ConversationId,
	})
	if err != nil {
		l.Errorf("[MarkRead] rpc err: %v convId: %d", err, req.ConversationId)
		return nil, err
	}

	return resp, nil
}
