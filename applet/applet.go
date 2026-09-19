package appletapi

import (
	"api-thinktalk/applet/internal/config"
	"api-thinktalk/applet/internal/handler"
	"api-thinktalk/applet/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type Config = config.Config

func RegisterRoutes(server *rest.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
}
