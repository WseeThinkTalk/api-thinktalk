package chat

import (
	"context"
	"encoding/json"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"
	user "api-thinktalk/client/user/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncOfflineLogic {
	return &SyncOfflineLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// unreadSyncItem 离线同步中的会话项
type unreadSyncItem struct {
	Id               int64  `json:"id"`
	TargetUserId     int64  `json:"target_user_id"`
	TargetUserName   string `json:"target_user_name"`
	TargetUserAvatar string `json:"target_user_avatar"`
	LastMessage      string `json:"last_message"`
	LastMessageTime  int64  `json:"last_message_time"`
	UnreadCount      int64  `json:"unread_count"`
}

// unreadSyncPayload 离线同步推送的 JSON 负载
type unreadSyncPayload struct {
	Type          string           `json:"type"`
	TotalUnread   int64            `json:"total_unread"`
	Conversations []unreadSyncItem `json:"conversations"`
}

// SyncOnConnect 用户 WebSocket 连接成功后调用，同步未读消息摘要
func (l *SyncOfflineLogic) SyncOnConnect(userId int64) {
	// 1. 查询最近会话（按最后消息时间倒序）
	resp, err := l.svcCtx.Chat.Conversations(l.ctx, &pb.ConversationsRequest{
		UserId:   userId,
		Cursor:   0,
		PageSize: 20,
	})
	if err != nil {
		l.Errorf("[SyncOffline] conversations err: %v userId: %d", err, userId)
		return
	}

	// 2. 只保留有未读消息的会话
	var unreadConvs []unreadSyncItem
	var totalUnread int64
	for _, item := range resp.Items {
		if item.UnreadCount == 0 {
			continue
		}
		totalUnread += item.UnreadCount

		conv := unreadSyncItem{
			Id:              item.Id,
			TargetUserId:    item.TargetUserId,
			LastMessage:     item.LastMessage,
			LastMessageTime: item.LastMessageTime,
			UnreadCount:     item.UnreadCount,
		}

		// 获取对方用户信息
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.TargetUserId}); err == nil {
			conv.TargetUserName = userResp.Username
			conv.TargetUserAvatar = userResp.Avatar
		}

		unreadConvs = append(unreadConvs, conv)
	}

	if totalUnread == 0 {
		return
	}

	// 3. 构建推送负载
	payload := unreadSyncPayload{
		Type:          "unread_sync",
		TotalUnread:   totalUnread,
		Conversations: unreadConvs,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		l.Errorf("[SyncOffline] marshal err: %v userId: %d", err, userId)
		return
	}

	// 4. 通过 WebSocket 推送给用户
	l.svcCtx.Hub.SendToUser(userId, &types.WsOutMessage{
		Type:    "unread_sync",
		Message: string(payloadBytes),
	})

	l.Infof("[SyncOffline] synced userId: %d unreadConvs: %d totalUnread: %d",
		userId, len(unreadConvs), totalUnread)
}
