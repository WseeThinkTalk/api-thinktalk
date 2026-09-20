package qa

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/qa"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func QuestionDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionDetailRequest
		if v := r.URL.Query().Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		l := logic.NewQuestionDetailLogic(r.Context(), svcCtx)
		resp, err := l.QuestionDetail(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
