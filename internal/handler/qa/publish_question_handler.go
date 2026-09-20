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

func PublishQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.PublishQuestionRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewPublishQuestionLogic(r.Context(), svcCtx)
		resp, err := l.PublishQuestion(uid, &req)
		errorx.HttpResult(r, w, resp, err)
	}
}
