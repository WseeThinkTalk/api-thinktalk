package social

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func FansListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.FansListRequest
		httpx.ParseQuery(r, &req)
		l := logic.NewFansListLogic(r.Context(), svcCtx)
		resp, err := l.FansList(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
