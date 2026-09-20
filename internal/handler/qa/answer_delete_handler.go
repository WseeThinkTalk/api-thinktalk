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

func AnswerDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := httpx.GetUserID(r)
		var req types.AnswerDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewAnswerDeleteLogic(r.Context(), svcCtx)
		err := l.AnswerDelete(uid, &req)
		errorx.HttpResult(r, w, nil, err)
	}
}
