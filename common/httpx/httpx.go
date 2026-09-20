package httpx

import (
	"encoding/json"
	"net/http"
)

func GetUserID(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(json.Number)
	uid, _ := userId.Int64()
	return uid
}

func WriteJSON(w http.ResponseWriter, data any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(data)
}

func ParseQuery(r *http.Request, v any) {
	if c := r.URL.Query().Get("cursor"); c != "" {
		json.Unmarshal([]byte(c), v)
	}
	if p := r.URL.Query().Get("page_size"); p != "" {
		json.Unmarshal([]byte(p), v)
	}
}
