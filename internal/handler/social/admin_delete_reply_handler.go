package social

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func AdminDeleteReplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDeleteReplyRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewAdminDeleteReplyLogic(r.Context(), svcCtx)
		resp, err := l.AdminDeleteReply(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
