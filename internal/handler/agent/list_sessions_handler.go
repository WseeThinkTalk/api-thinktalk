package agent

import (
	"net/http"

	logic "api-thinktalk/internal/logic/agent"
	"api-thinktalk/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListSessionsLogic(r.Context(), svcCtx)
		resp, err := l.ListSessions()
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}
