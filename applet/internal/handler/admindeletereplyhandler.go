package handler

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/applet/internal/logic"
	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
)

func AdminDeleteReplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDeleteReplyRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewAdminDeleteReplyLogic(r.Context(), svcCtx)
		resp, err := l.AdminDeleteReply(&req)
		writeJSON(w, resp, err)
	}
}
