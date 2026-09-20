package tag

import (
	"api-thinktalk/common/httpx"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/tag"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func TagListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagListRequest
		httpx.ParseQuery(r, &req)
		l := logic.NewTagListLogic(r.Context(), svcCtx)
		resp, err := l.TagList(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
