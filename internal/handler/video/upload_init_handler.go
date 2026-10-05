package video

import (
	"net/http"

	logic "api-thinktalk/internal/logic/video"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/lib/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadInitHandler 视频分片上传初始化 Handler
func UploadInitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoUploadInitRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadInitLogic(r.Context(), svcCtx)
		resp, err := l.UploadInit(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
