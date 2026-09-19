package svc

import (
	"api-thinktalk/agent/internal/config"
	"api-thinktalk/client/agent/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	AgentClient pb.AgentClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: c.AgentRpc.Endpoints,
		NonBlock:  c.AgentRpc.NonBlock,
		Timeout:   c.AgentRpc.Timeout,
	}).Conn()

	return &ServiceContext{
		Config:      c,
		AgentClient: pb.NewAgentClient(conn),
	}
}
