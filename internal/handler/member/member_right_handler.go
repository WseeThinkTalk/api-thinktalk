package member

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/member"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func MemberRightHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.MemberRightRequest
		if k := r.URL.Query().Get("right_key"); k != "" {
			req.RightKey = k
		}
		l := logic.NewMemberRightLogic(r.Context(), svcCtx)
		resp, err := l.CheckRight(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
