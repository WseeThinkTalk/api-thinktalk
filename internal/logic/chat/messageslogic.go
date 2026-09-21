package chat

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagesLogic {
	return &MessagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MessagesLogic) Messages(userId int64, req *types.MessagesRequest) (resp *types.MessagesResponse, err error) {
	resp = new(types.MessagesResponse)

	if req.PageSize == 0 {
		req.PageSize = types.DefaultPageSize
	}

	rpcResp, err := l.svcCtx.Chat.Messages(l.ctx, &pb.MessagesRequest{
		ConversationId: req.ConversationId,
		Cursor:         req.Cursor,
		PageSize:       req.PageSize,
	})
	if err != nil {
		l.Errorf("[Messages] rpc err: %v convId: %d", err, req.ConversationId)
		return nil, err
	}

	// 转换聊天消息数据项
	items := make([]*types.MessageItem, 0, len(rpcResp.Items))
	for _, v := range rpcResp.Items {
		items = append(items, &types.MessageItem{
			Id:             v.Id,
			ConversationId: v.ConversationId,
			SenderId:       v.SenderId,
			ReceiverId:     v.ReceiverId,
			Content:        v.Content,
			MsgType:        v.MsgType,
			IsRead:         v.IsRead,
			CreateTime:     v.CreateTime,
		})
	}

	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
