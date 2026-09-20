package social

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
)

func UnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		l := logic.NewUnreadCountLogic(r.Context(), svcCtx)
		resp, err := l.UnreadCount(uid)
		errorx.HttpResult(r, w, resp, err)
	}
}
