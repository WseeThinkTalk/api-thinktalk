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

func TagCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.TagCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagCreateLogic(r.Context(), svcCtx)
		resp, err := l.TagCreate(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
