package social

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ReplyDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyDetailRequest
		if v := r.URL.Query().Get("reply_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ReplyId)
		}
		l := logic.NewReplyDetailLogic(r.Context(), svcCtx)
		resp, err := l.ReplyDetail(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
