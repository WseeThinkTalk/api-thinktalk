package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/agent/pb"
)

type ListSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionsLogic {
	return &ListSessionsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListSessionsLogic) ListSessions() (*types.ListSessionsResponse, error) {
	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.AgentClient.ListSessions(l.ctx, &pb.ListSessionsRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.SessionItem, len(resp.Sessions))
	for i, s := range resp.Sessions {
		items[i] = types.SessionItem{
			SessionID:    s.SessionId,
			Title:        s.Title,
			MessageCount: s.MessageCount,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		}
	}
	return &types.ListSessionsResponse{Sessions: items}, nil
}
