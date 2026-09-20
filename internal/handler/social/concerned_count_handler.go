package social

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ConcernedCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConcernedCountRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("obj_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ObjId)
		}
		l := logic.NewConcernedCountLogic(r.Context(), svcCtx)
		resp, err := l.ConcernedCount(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
