package chatapi

import (
	"api-thinktalk/chat/internal/config"
	"api-thinktalk/chat/internal/handler"
	"api-thinktalk/chat/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type Config = config.Config

func RegisterRoutes(server *rest.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
}
