package articleapi

import (
	"api-thinktalk/article/internal/config"
	"api-thinktalk/article/internal/handler"
	"api-thinktalk/article/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type Config = config.Config

func RegisterRoutes(server *rest.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
}
