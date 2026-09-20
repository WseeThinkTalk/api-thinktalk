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

func ConcernedCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.ConcernedCheckRequest
		if b := r.URL.Query().Get("biz_id"); b != "" {
			req.BizId = b
		}
		if o := r.URL.Query().Get("obj_id"); o != "" {
			json.Unmarshal([]byte(o), &req.ObjId)
		}
		l := logic.NewConcernedCheckLogic(r.Context(), svcCtx)
		resp, err := l.Check(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
