package tag

import (
	"api-thinktalk/common/httpx"
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func UntagResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.UntagResourceRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewUntagResourceLogic(r.Context(), svcCtx)
		resp, err := l.UntagResource(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
