package handler

import (
	"net/http"

	"api-thinktalk/chat/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/conversations",
				Handler: ConversationsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/messages",
				Handler: MessagesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/send",
				Handler: SendMessageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/markread",
				Handler: MarkReadHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/unread",
				Handler: UnreadCountHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/online",
				Handler: OnlineHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/chat"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ws",
				Handler: WebSocketHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1/chat"),
	)
}
