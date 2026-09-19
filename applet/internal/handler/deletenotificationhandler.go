package handler

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/applet/internal/logic"
	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
)

func DeleteNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.DeleteNotificationRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewDeleteNotificationLogic(r.Context(), svcCtx)
		resp, err := l.DeleteNotification(uid, &req)
		writeJSON(w, resp, err)
	}
}
