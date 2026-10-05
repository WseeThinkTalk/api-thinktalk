package video

import (
	"net/http"

	logic "api-thinktalk/internal/logic/video"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/lib/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// VideoStatusHandler 视频状态查询 Handler
func VideoStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewVideoStatusLogic(r.Context(), svcCtx)
		resp, err := l.VideoStatus(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
