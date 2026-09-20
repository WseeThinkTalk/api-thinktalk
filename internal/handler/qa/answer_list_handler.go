package qa

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/qa"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func AnswerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AnswerListRequest
		q := r.URL.Query()
		if v := q.Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewAnswerListLogic(r.Context(), svcCtx)
		resp, err := l.AnswerList(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
