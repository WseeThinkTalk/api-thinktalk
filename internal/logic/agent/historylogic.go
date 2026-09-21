package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/agent/pb"
)

type HistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryLogic {
	return &HistoryLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *HistoryLogic) History(req *types.HistoryRequest) (resp *types.HistoryResponse, err error) {
	resp = new(types.HistoryResponse)

	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AgentClient.GetHistory(l.ctx, &pb.GetHistoryRequest{
		UserId:    uid,
		SessionId: req.SessionID,
	})
	if err != nil {
		return nil, err
	}

	// 转换会话历史消息列表
	msgs := make([]types.AgentMessageItem, len(rpcResp.Messages))
	for i, v := range rpcResp.Messages {
		msgs[i] = types.AgentMessageItem{Role: v.Role, Content: v.Content}
	}
	resp.SessionID = rpcResp.SessionId
	resp.Title = rpcResp.Title
	resp.Messages = msgs
	resp.CreatedAt = rpcResp.CreatedAt
	resp.UpdatedAt = rpcResp.UpdatedAt
	return resp, nil
}
