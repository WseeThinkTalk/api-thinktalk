package social

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	user "api-thinktalk/client/user/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type AdminBroadcastLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminBroadcastLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminBroadcastLogic {
	return &AdminBroadcastLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminBroadcastLogic) AdminBroadcast(req *types.AdminBroadcastNotificationRequest) (*types.AdminBroadcastNotificationResponse, error) {
	if req.Title == "" || req.Content == "" {
		return nil, fmt.Errorf("标题和内容不能为空")
	}

	l.Infof("[AdminBroadcast] broadcasting system notification title: %s", req.Title)

	var totalCount int64
	var successCount int64

	// 分页获取所有用户ID
	cursor := int64(0)
	pageSize := int64(200)

	for {
		rpcResp, err := l.svcCtx.UserRPC.AdminUserList(l.ctx, &user.AdminUserListRequest{
			Keyword:  "",
			Cursor:   cursor,
			PageSize: pageSize,
		})
		if err != nil {
			l.Errorf("[AdminBroadcast] AdminUserList err: %v cursor: %d", err, cursor)
			return nil, err
		}

		if len(rpcResp.Items) == 0 {
			break
		}

		totalCount += int64(len(rpcResp.Items))

		// 异步推送通知到每个用户
		title := req.Title
		content := req.Content
		pusher := l.svcCtx.NotificationPusher
		for _, item := range rpcResp.Items {
			userId := item.UserId
			threading.GoSafe(func() {
				notif := map[string]interface{}{
					"userId":        userId,
					"type":          int32(4), // NotifTypeSystem
					"title":         title,
					"content":       content,
					"refId":         int64(0),
					"bizId":         fmt.Sprintf("system_broadcast_%d", userId),
					"triggerUserId": int64(0), // 系统通知无触发用户
				}
				data, err := json.Marshal(notif)
				if err != nil {
					l.Errorf("[AdminBroadcast] marshal err: %v", err)
					return
				}
				if err := pusher.Push(context.Background(), string(data)); err != nil {
					l.Errorf("[AdminBroadcast] push err for userId %d: %v", userId, err)
				}
			})
		}

		successCount += int64(len(rpcResp.Items))

		if rpcResp.IsEnd {
			break
		}
		cursor = rpcResp.Cursor
	}

	l.Infof("[AdminBroadcast] done total: %d success: %d", totalCount, successCount)
	return &types.AdminBroadcastNotificationResponse{
		TotalCount:   totalCount,
		SuccessCount: successCount,
	}, nil
}
