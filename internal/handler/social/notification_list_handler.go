package social

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"
	"strconv"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func NotificationListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.NotificationRequest
		if t := r.URL.Query().Get("type"); t != "" {
			v, _ := strconv.Atoi(t)
			req.Type = int32(v)
		}
		if c := r.URL.Query().Get("cursor"); c != "" {
			v, _ := strconv.ParseInt(c, 10, 64)
			req.Cursor = v
		}
		if p := r.URL.Query().Get("page_size"); p != "" {
			v, _ := strconv.ParseInt(p, 10, 64)
			req.PageSize = v
		}
		l := logic.NewNotificationListLogic(r.Context(), svcCtx)
		resp, err := l.NotificationList(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
