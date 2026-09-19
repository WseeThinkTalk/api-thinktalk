package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-thinktalk/applet/internal/logic"
	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
)

func NotificationListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.NotificationRequest
		if t := r.URL.Query().Get("type"); t != "" {
			v, _ := strconv.Atoi(t)
			req.Type = int32(v)
		}
		if c := r.URL.Query().Get("cursor"); c != "" {
			v, _ := strconv.ParseInt(c, 10, 64)
			req.Cursor = v
		}
		if p := r.URL.Query().Get("page_size"); p != "" {
			v, _ := strconv.ParseInt(p, 10, 64)
			req.PageSize = v
		}
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.NotificationList(uid, &req)
		writeJSON(w, resp, err)
	}
}

func UnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.UnreadCount(uid)
		writeJSON(w, resp, err)
	}
}

func MarkReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MarkReadRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.MarkRead(uid, &req)
		writeJSON(w, resp, err)
	}
}

func MarkAllReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MarkAllReadRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.MarkAllRead(uid, &req)
		writeJSON(w, resp, err)
	}
}
