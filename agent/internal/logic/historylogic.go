package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/agent/internal/svc"
	"api-thinktalk/agent/internal/types"
	"api-thinktalk/client/agent/pb"
)

type HistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryLogic {
	return &HistoryLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *HistoryLogic) History(req *types.HistoryRequest) (*types.HistoryResponse, error) {
	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.AgentClient.GetHistory(l.ctx, &pb.GetHistoryRequest{
		UserId:    uid,
		SessionId: req.SessionID,
	})
	if err != nil {
		return nil, err
	}

	msgs := make([]types.MessageItem, len(resp.Messages))
	for i, m := range resp.Messages {
		msgs[i] = types.MessageItem{Role: m.Role, Content: m.Content}
	}
	return &types.HistoryResponse{
		SessionID: resp.SessionId,
		Title:     resp.Title,
		Messages:  msgs,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}
