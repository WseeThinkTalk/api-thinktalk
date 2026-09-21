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

	if req.PageSize == 0 {
		req.PageSize = types.DefaultPageSize
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

	items := make([]*types.ConversationItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		targetName := ""
		targetAvatar := ""
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.TargetUserId}); err == nil {
			targetName = userResp.Username
			targetAvatar = userResp.Avatar
		}

		items = append(items, &types.ConversationItem{
			Id:               item.Id,
			TargetUserId:     item.TargetUserId,
			TargetUserName:   targetName,
			TargetUserAvatar: targetAvatar,
			LastMessage:      item.LastMessage,
			LastMessageTime:  item.LastMessageTime,
			UnreadCount:      item.UnreadCount,
		})
	}

	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
