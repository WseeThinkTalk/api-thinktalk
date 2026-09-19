package agentapi

import (
	"api-thinktalk/agent/internal/config"
	"api-thinktalk/agent/internal/handler"
	"api-thinktalk/agent/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type Config = config.Config

func RegisterRoutes(server *rest.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
}
