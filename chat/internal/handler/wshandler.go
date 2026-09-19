package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"api-thinktalk/chat/internal/logic"
	"api-thinktalk/chat/internal/svc"
	"api-thinktalk/chat/internal/types"
	"api-thinktalk/client/chat/pb"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var allowedOrigins = map[string]bool{
	"":                          true, // no Origin header (native apps, Postman)
	"http://localhost:5173":     true, // dev frontend
	"http://localhost:3000":     true, // prod frontend
	"http://127.0.0.1:5173":    true,
	"http://127.0.0.1:3000":    true,
}

func isOriginAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		return true
	}
	logx.Infof("[WS] rejected origin: %s", origin)
	return false
}

var upgrader = websocket.Upgrader{
	CheckOrigin: isOriginAllowed,
}

func WebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	hub := svcCtx.Hub

	hub.OnMessage = func(userId int64, msg *types.WsInMessage) {
		switch msg.Type {
		case "message":
			_, err := svcCtx.Chat.SendMessage(context.Background(), &pb.SendMessageRequest{
				SenderId:   userId,
				ReceiverId: msg.ReceiverId,
				Content:    msg.Content,
				MsgType:    msg.MsgType,
			})
			if err != nil {
				logx.Errorf("[WS] send message err: %v userId: %d", err, userId)
				hub.SendToUser(userId, &types.WsOutMessage{
					Type:    "error",
					Message: err.Error(),
				})
				return
			}
			// 即时推送给接收者（在线则实时收到，离线则静默跳过）
			hub.SendToUser(msg.ReceiverId, &types.WsOutMessage{
				Type:        "new_message",
				SenderId:    userId,
				Content:     msg.Content,
				MsgType:     msg.MsgType,
				ClientMsgId: msg.ClientMsgId,
			})
			hub.SendToUser(userId, &types.WsOutMessage{
				Type:    "ack",
				Message: "sent",
			})
		default:
			hub.SendToUser(userId, &types.WsOutMessage{
				Type:    "error",
				Message: "unknown message type",
			})
		}
	}

	hub.OnConnect = func(userId int64) {
		l := logic.NewSyncOfflineLogic(context.Background(), svcCtx)
		l.SyncOnConnect(userId)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := ""
		// 优先从 Authorization Header 取
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
			tokenStr = authHeader[7:]
		}
		// 浏览器 WebSocket 无法自定义 Header，fallback 到 query string
		if tokenStr == "" {
			tokenStr = r.URL.Query().Get("token")
		}

		if tokenStr == "" {
			http.Error(w, "unauthorized: token missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(svcCtx.Config.Auth.AccessSecret), nil
		})
		if err != nil || !token.Valid {
			logx.Errorf("[WS] token parse err: %v", err)
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized: invalid claims", http.StatusUnauthorized)
			return
		}

		var uid int64
		if userIdVal, exists := claims["userId"]; exists {
			switch val := userIdVal.(type) {
			case float64:
				uid = int64(val)
			case json.Number:
				uid, _ = val.Int64()
			}
		}

		if uid == 0 {
			http.Error(w, "unauthorized: invalid user id", http.StatusUnauthorized)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("[WS] upgrade err: %v", err)
			return
		}

		conn := hub.Register(uid, ws)
		logx.Infof("[WS] user connected userId: %d", uid)

		// 用户上线后触发离线消息同步
		if hub.OnConnect != nil {
			hub.OnConnect(uid)
		}

		go hub.WritePump(conn)
		go hub.ReadPump(conn)
	}
}
