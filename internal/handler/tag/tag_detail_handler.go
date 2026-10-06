package tag

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func TagDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagDetailRequest
		if v := r.URL.Query().Get("id"); v != "" {
			json.Unmarshal([]byte(v), &req.Id)
		} else if v := r.URL.Query().Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.Id)
		}
		l := logic.NewTagDetailLogic(r.Context(), svcCtx)
		resp, err := l.TagDetail(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
