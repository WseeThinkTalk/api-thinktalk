package video

import (
	"net/http"

	logic "api-thinktalk/internal/logic/video"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/lib/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadCompleteHandler 视频分片合并 Handler
func UploadCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoUploadCompleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadCompleteLogic(r.Context(), svcCtx)
		resp, err := l.UploadComplete(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
