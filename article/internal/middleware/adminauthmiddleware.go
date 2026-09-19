package middleware

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/client/user/user"
)

type AdminAuthMiddleware struct {
	UserRPC user.User
}

func NewAdminAuthMiddleware(userRPC user.User) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		UserRPC: userRPC,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdVal := r.Context().Value("userId")
		if userIdVal == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		num, ok := userIdVal.(json.Number)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, err := num.Int64()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		u, err := m.UserRPC.FindById(r.Context(), &user.FindByIdRequest{UserId: userId})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if u.Role != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"code":403,"msg":"权限不足，需要管理员权限"}`))
			return
		}

		next(w, r)
	}
}
