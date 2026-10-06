package social

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	msg "api-thinktalk/client/message/pb"

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

func (l *DeleteNotificationLogic) DeleteNotification(userId int64, req *types.DeleteNotificationRequest) (resp *types.DeleteNotificationResponse, err error) {
	resp = new(types.DeleteNotificationResponse)

	_, err = l.svcCtx.MessageRPC.DeleteNotification(l.ctx, &msg.DeleteNotificationRequest{
		UserId:         userId,
		NotificationId: req.Id,
	})
	if err != nil {
		l.Errorf("[DeleteNotification] rpc err: %v", err)
		return nil, err
	}
	return resp, nil
}
