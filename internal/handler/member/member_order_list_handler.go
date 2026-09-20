package member

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/member"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func MemberOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.MemberOrderListRequest
		httpx.ParseQuery(r, &req)
		l := logic.NewMemberOrderListLogic(r.Context(), svcCtx)
		resp, err := l.MemberOrderList(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
