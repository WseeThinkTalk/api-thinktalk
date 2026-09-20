package qa

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	"api-thinktalk/common/httpx"
	logic "api-thinktalk/internal/logic/qa"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func QuestionDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.QuestionDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQuestionDeleteLogic(r.Context(), svcCtx)
		err := l.QuestionDelete(uid, &req)
		errorx.HttpResult(r, w, nil, err)
	}
}
