package svc

import (
	"api-thinktalk/qa/internal/config"
	"api-thinktalk/client/qa/qa"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	QaRPC  qa.QA
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		QaRPC:  qa.NewQA(zrpc.MustNewClient(c.QaRpc)),
	}
}
