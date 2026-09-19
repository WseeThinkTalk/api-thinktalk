package logic

import (
	"context"

	"api-thinktalk/chat/internal/svc"
	"api-thinktalk/chat/internal/types"
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

func (l *UnreadCountLogic) UnreadCount(userId int64) (*types.UnreadCountResponse, error) {
	resp, err := l.svcCtx.Chat.UnreadCount(l.ctx, &pb.UnreadCountRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[UnreadCount] rpc err: %v userId: %d", err, userId)
		return nil, err
	}

	return &types.UnreadCountResponse{Total: resp.Total}, nil
}
