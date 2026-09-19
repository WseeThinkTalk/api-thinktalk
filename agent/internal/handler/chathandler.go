package handler

import (
	"net/http"

	"api-thinktalk/agent/internal/logic"
	"api-thinktalk/agent/internal/svc"
	"api-thinktalk/agent/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewChatLogic(r.Context(), svcCtx)
		l.Chat(&req, w)
	}
}
