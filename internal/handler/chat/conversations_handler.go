package chat

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/chat"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req types.ConversationsRequest
		if v := r.URL.Query().Get("cursor"); v != "" {
			_ = json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := r.URL.Query().Get("page_size"); v != "" {
			_ = json.Unmarshal([]byte(v), &req.PageSize)
		}

		l := logic.NewConversationsLogic(r.Context(), svcCtx)
		resp, err := l.Conversations(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
