package social

import (
	"api-thinktalk/common/httpx"
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ConcernedCancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.ConcernedCancelRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewConcernedCancelLogic(r.Context(), svcCtx)
		resp, err := l.Cancel(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
