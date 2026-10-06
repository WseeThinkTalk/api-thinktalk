package tag

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func TagsByResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagsByResourceRequest
		q := r.URL.Query()
		if v := q.Get("resource_type"); v != "" {
			req.ResourceType = v
		} else if v := q.Get("biz_id"); v != "" {
			req.ResourceType = v
		}
		if v := q.Get("resource_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ResourceId)
		} else if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ResourceId)
		}
		l := logic.NewTagsByResourceLogic(r.Context(), svcCtx)
		resp, err := l.TagsByResource(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
