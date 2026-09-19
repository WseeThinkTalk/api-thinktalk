package handler

import (
	"net/http"

	"api-thinktalk/applet/internal/logic"
	"api-thinktalk/applet/internal/svc"
)

func AlipayNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAlipayNotifyLogic(r.Context(), svcCtx)
		err := l.AlipayNotify(r)
		if err != nil {
			// 如果处理失败或验签失败，支付宝要求返回 "failure" （或非 "success"）
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failure"))
		} else {
			// 处理成功，返回 success，让支付宝停止重试
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}
}
