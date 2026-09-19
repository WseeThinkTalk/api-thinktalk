package qaapi

import (
	"api-thinktalk/qa/internal/config"
	"api-thinktalk/qa/internal/handler"
	"api-thinktalk/qa/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type Config = config.Config

func RegisterRoutes(server *rest.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
}
