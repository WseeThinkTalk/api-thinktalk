package social

import (
	"encoding/json"
	"net/http"

	"api-thinktalk/pkg/lib/errorx"

	logic "api-thinktalk/internal/logic/social"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
)

func ReplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyListRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("obj_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ObjId)
		} else if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ObjId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewReplyListLogic(r.Context(), svcCtx)
		resp, err := l.ReplyList(&req)
		errorx.HttpResult(r, w, resp, err)
	}
}
