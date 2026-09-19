package logic

import (
	"context"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	msg "api-thinktalk/client/message/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationLogic {
	return &DeleteNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNotificationLogic) DeleteNotification(userId int64, req *types.DeleteNotificationRequest) (*types.DeleteNotificationResponse, error) {
	_, err := l.svcCtx.MessageRPC.DeleteNotification(l.ctx, &msg.DeleteNotificationRequest{
		UserId:         userId,
		NotificationId: req.NotificationId,
	})
	if err != nil {
		l.Errorf("[DeleteNotification] rpc err: %v", err)
		return nil, err
	}
	return &types.DeleteNotificationResponse{}, nil
}
