package handler

import (
	"net/http"

	"api-thinktalk/applet/internal/logic"
	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminBroadcastHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminBroadcastNotificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAdminBroadcastLogic(r.Context(), svcCtx)
		resp, err := l.AdminBroadcast(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
