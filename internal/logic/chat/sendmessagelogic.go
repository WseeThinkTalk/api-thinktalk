package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(userId int64, req *types.SendMessageRequest) (*types.SendMessageResponse, error) {
	_, err := l.svcCtx.Chat.SendMessage(l.ctx, &pb.SendMessageRequest{
		SenderId:   userId,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
		MsgType:    req.MsgType,
	})
	if err != nil {
		l.Errorf("[SendMessage] rpc err: %v req: %+v", err, req)
		return nil, err
	}

	// 即时推送给接收者（在线则实时收到，离线则消息已由 Kafka→DB 持久化）
	l.svcCtx.Hub.SendToUser(req.ReceiverId, &types.WsOutMessage{
		Type:     "new_message",
		SenderId: userId,
		Content:  req.Content,
		MsgType:  req.MsgType,
	})

	return &types.SendMessageResponse{}, nil
}
