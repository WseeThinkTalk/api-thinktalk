package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/agent/pb"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionRequest) (resp *types.DeleteSessionResponse, err error) {
	resp = new(types.DeleteSessionResponse)

	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AgentClient.DeleteSession(l.ctx, &pb.DeleteSessionRequest{
		UserId:    uid,
		SessionId: req.SessionID,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp != nil && rpcResp.Data != nil {
		resp.Success = rpcResp.Data.Success
	}
	return resp, nil
}
