package article

import (
	"net/http"

	logic "api-thinktalk/internal/logic/article"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/lib/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ArticleDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewArticleDeleteLogic(r.Context(), svcCtx)
		err := l.ArticleDelete(&req)
		errorx.HttpResult(r, w, nil, err)
	}
}
