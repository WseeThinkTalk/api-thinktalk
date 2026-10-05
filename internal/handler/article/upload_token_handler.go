package article

import (
	"net/http"

	logic "api-thinktalk/internal/logic/article"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/lib/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadTokenHandler 获取上传凭证
func UploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadTokenLogic(r.Context(), svcCtx)
		resp, err := l.UploadToken(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
