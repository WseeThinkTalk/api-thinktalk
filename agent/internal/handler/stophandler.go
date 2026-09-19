package handler

import (
	"net/http"

	"api-thinktalk/agent/internal/logic"
	"api-thinktalk/agent/internal/svc"
	"api-thinktalk/agent/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func StopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StopRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewStopLogic(r.Context(), svcCtx)
		if err := l.Stop(&req); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.Ok(w)
	}
}
