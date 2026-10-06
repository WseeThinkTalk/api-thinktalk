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

func (l *NotificationListLogic) NotificationList(userId int64, req *types.NotificationRequest) (resp *types.NotificationResponse, err error) {
	resp = new(types.NotificationResponse)
	resp.Items = make([]*types.NotificationItem, 0)

	rpcResp, err := l.svcCtx.MessageRPC.NotificationList(l.ctx, &pb.NotificationListRequest{
		UserId:   userId,
		Type:     req.Type,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[NotificationList] rpc err: %v", err)
		return nil, err
	}

	// 转换通知列表项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.NotificationItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			items = append(items, &types.NotificationItem{
				Id:         v.Id,
				UserId:     userId,
				SenderId:   v.TriggerUserId,
				MsgType:    v.Type,
				Title:      v.Title,
				Content:    v.Content,
				BizId:      v.BizId,
				ObjId:      v.RefId,
				IsRead:     v.IsRead,
				CreateTime: v.CreateTime,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
