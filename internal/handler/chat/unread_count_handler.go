package chat

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/chat"
	"api-thinktalk/internal/svc"
)

func UnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		l := logic.NewUnreadCountLogic(r.Context(), svcCtx)
		resp, err := l.UnreadCount(uid)
		errorx.HttpResult(r, w, resp, err)
	}
}
