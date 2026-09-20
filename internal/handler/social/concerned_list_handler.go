package social

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ConcernedListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.ConcernedListRequest
		if b := r.URL.Query().Get("biz_id"); b != "" {
			req.BizId = b
		}
		httpx.ParseQuery(r, &req)
		l := logic.NewConcernedListLogic(r.Context(), svcCtx)
		resp, err := l.List(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
