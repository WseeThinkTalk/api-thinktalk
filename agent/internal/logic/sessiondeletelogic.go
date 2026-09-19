package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/agent/internal/svc"
	"api-thinktalk/agent/internal/types"
	"api-thinktalk/client/agent/pb"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionRequest) (*types.DeleteSessionResponse, error) {
	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.AgentClient.DeleteSession(l.ctx, &pb.DeleteSessionRequest{
		UserId:    uid,
		SessionId: req.SessionID,
	})
	if err != nil {
		return &types.DeleteSessionResponse{Success: false}, err
	}
	return &types.DeleteSessionResponse{Success: resp.Success}, nil
}
