package social

import (
	"context"

	"api-thinktalk/client/message/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotificationListLogic) NotificationList(userId int64, req *types.NotificationRequest) (*types.NotificationResponse, error) {
	resp, err := l.svcCtx.MessageRPC.NotificationList(l.ctx, &pb.NotificationListRequest{
		UserId:   userId,
		Type:     req.Type,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[NotificationList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.NotificationItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.NotificationItem{
			Id:            item.Id,
			Type:          item.Type,
			Title:         item.Title,
			Content:       item.Content,
			IsRead:        item.IsRead,
			TriggerUserId: item.TriggerUserId,
			RefId:         item.RefId,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.NotificationResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
