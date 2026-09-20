package chat

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/chat"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func MarkReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req types.ChatMarkReadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := logic.NewMarkReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkRead(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
