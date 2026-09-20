package member

import (
	"api-thinktalk/common/httpx"
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/member"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func UpgradeMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.UpgradeMemberRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewUpgradeMemberLogic(r.Context(), svcCtx)
		resp, err := l.UpgradeMember(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
