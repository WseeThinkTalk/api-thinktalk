package social

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func FollowListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.FollowListRequest
		httpx.ParseQuery(r, &req)
		l := logic.NewFollowListLogic(r.Context(), svcCtx)
		resp, err := l.FollowList(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
