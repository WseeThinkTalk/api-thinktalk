package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"
	user "api-thinktalk/client/user/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationsLogic {
	return &ConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConversationsLogic) Conversations(userId int64, req *types.ConversationsRequest) (resp *types.ConversationsResponse, err error) {
	resp = new(types.ConversationsResponse)

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	rpcResp, err := l.svcCtx.Chat.Conversations(l.ctx, &pb.ConversationsRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[Conversations] rpc err: %v userId: %d", err, userId)
		return nil, err
	}

	// 转换会话列表并补充对方用户信息
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.ConversationItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			targetName := ""
			targetAvatar := ""
			if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: v.TargetUserId}); err == nil && userResp != nil && userResp.Data != nil {
				targetName = userResp.Data.Username
				targetAvatar = userResp.Data.Avatar
			}

			items = append(items, &types.ConversationItem{
				Id:               v.Id,
				TargetUserId:     v.TargetUserId,
				TargetUserName:   targetName,
				TargetUserAvatar: targetAvatar,
				LastMessage:      v.LastMessage,
				LastMessageTime:  v.LastMessageTime,
				UnreadCount:      v.UnreadCount,
			})
		}

		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
