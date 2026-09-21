package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type WsMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWsMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsMessageLogic {
	return &WsMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *WsMessageLogic) HandleMessage(userId int64, msg *types.WsInMessage) (resp *types.WsOutMessage, err error) {
	resp = new(types.WsOutMessage)

	_, err = l.svcCtx.Chat.SendMessage(l.ctx, &pb.SendMessageRequest{
		SenderId:   userId,
		ReceiverId: msg.ReceiverId,
		Content:    msg.Content,
		MsgType:    msg.MsgType,
	})
	if err != nil {
		l.Errorf("[WsMessage] rpc err: %v userId: %d", err, userId)
		return nil, err
	}

	resp.Type = "ack"
	resp.Message = "sent"
	return resp, nil
}
