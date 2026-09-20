package tag

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func HotTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotTagsRequest
		if v := r.URL.Query().Get("limit"); v != "" {
			json.Unmarshal([]byte(v), &req.Limit)
		}
		l := logic.NewHotTagsLogic(r.Context(), svcCtx)
		resp, err := l.HotTags(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
