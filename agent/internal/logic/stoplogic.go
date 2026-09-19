package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/agent/internal/svc"
	"api-thinktalk/agent/internal/types"
	"api-thinktalk/client/agent/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type StopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopLogic {
	return &StopLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *StopLogic) Stop(req *types.StopRequest) error {
	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		return fmt.Errorf("unauthorized")
	}
	uid, err := userIDVal.(json.Number).Int64()
	if err != nil {
		return err
	}

	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: l.svcCtx.Config.AgentRpc.Endpoints,
		NonBlock:  l.svcCtx.Config.AgentRpc.NonBlock,
		Timeout:   l.svcCtx.Config.AgentRpc.Timeout,
	})

	conn := client.Conn()
	_, err = pb.NewAgentClient(conn).Stop(l.ctx, &pb.StopRequest{
		UserId:    uid,
		SessionId: req.SessionID,
	})
	return err
}
