package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-thinktalk/chat/internal/svc"
)

func OnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		targetUserStr := r.URL.Query().Get("userId")
		if targetUserStr == "" {
			http.Error(w, "missing userId query parameter", http.StatusBadRequest)
			return
		}

		targetUserId, err := strconv.ParseInt(targetUserStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid userId query parameter", http.StatusBadRequest)
			return
		}

		isOnline := svcCtx.Hub.IsOnline(targetUserId)

		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"is_online": isOnline,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
