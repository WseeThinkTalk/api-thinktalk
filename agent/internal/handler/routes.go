package handler

import (
	"net/http"

	"api-thinktalk/agent/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodPost, Path: "/chat", Handler: ChatHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/stop", Handler: StopHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/history", Handler: HistoryHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/sessions", Handler: ListSessionsHandler(serverCtx)},
			{Method: http.MethodDelete, Path: "/sessions", Handler: DeleteSessionHandler(serverCtx)},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/agent"),
	)
}
