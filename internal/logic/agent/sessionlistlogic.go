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

func (l *ListSessionsLogic) ListSessions() (resp *types.ListSessionsResponse, err error) {
	resp = new(types.ListSessionsResponse)

	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AgentClient.ListSessions(l.ctx, &pb.ListSessionsRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	// 组装智能体历史会话项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]types.SessionItem, len(rpcResp.Data.Sessions))
		for i, v := range rpcResp.Data.Sessions {
			items[i] = types.SessionItem{
				SessionID:    v.SessionId,
				Title:        v.Title,
				MessageCount: v.MessageCount,
				CreatedAt:    v.CreatedAt,
				UpdatedAt:    v.UpdatedAt,
			}
		}
		resp.Sessions = items
	}
	return resp, nil
}
