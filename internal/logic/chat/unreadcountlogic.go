package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadCountLogic {
	return &UnreadCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnreadCountLogic) UnreadCount(userId int64) (resp *types.ChatUnreadCountResponse, err error) {
	resp = new(types.ChatUnreadCountResponse)

	rpcResp, err := l.svcCtx.Chat.UnreadCount(l.ctx, &pb.ChatUnreadCountRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[UnreadCount] rpc err: %v userId: %d", err, userId)
		return nil, err
	}

	if rpcResp != nil && rpcResp.Data != nil {
		resp.Total = rpcResp.Data.Total
	}
	return resp, nil
}
