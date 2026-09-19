// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"api-thinktalk/article/internal/logic"
	"api-thinktalk/article/internal/svc"
	"api-thinktalk/article/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AllArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AllArticlesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAllArticlesLogic(r.Context(), svcCtx)
		resp, err := l.AllArticles(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
