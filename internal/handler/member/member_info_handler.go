package member

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/member"
	"api-thinktalk/internal/svc"
)

func MemberInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		l := logic.NewMemberInfoLogic(r.Context(), svcCtx)
		resp, err := l.MemberInfo(uid)
		errorx.HttpResult(r, w, resp, err)
	}
}
