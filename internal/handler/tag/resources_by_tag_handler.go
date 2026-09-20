package tag

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ResourcesByTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResourcesByTagRequest
		q := r.URL.Query()
		if v := q.Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TagId)
		}
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewResourcesByTagLogic(r.Context(), svcCtx)
		resp, err := l.ResourcesByTag(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
