package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/agent/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

type ChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger logx.Logger
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{ctx: ctx, svcCtx: svcCtx, logger: logx.WithContext(ctx)}
}

func (l *ChatLogic) Chat(req *types.ChatRequest, w http.ResponseWriter) {
	userIDVal := l.ctx.Value("userId")
	if userIDVal == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := userIDVal.(json.Number).Int64()
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("sess_%d_%d", userID, time.Now().UnixMilli())
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Session-Id", sessionID)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	stream, err := l.svcCtx.AgentClient.Chat(l.ctx, &pb.ChatRequest{
		UserId:    userID,
		SessionId: sessionID,
		Message:   req.Message,
	})
	if err != nil {
		l.writeSSE(w, flusher, "error", err.Error(), "", "")
		return
	}

	for {
		evt, err := stream.Recv()
		if err != nil {
			l.logger.Errorf("[agent-api] Recv err: %v", err)
			l.writeSSE(w, flusher, "error", err.Error(), "", "")
			break
		}
		if evt.Type == "" && evt.Content == "" {
			continue
		}
		if evt.Type == "done" {
			l.writeSSE(w, flusher, "done", evt.FinishReason, "", "")
			break
		}
		l.writeSSE(w, flusher, evt.Type, evt.Content, evt.ToolName, evt.ToolData)
	}
	l.writeSSE(w, flusher, "done", "stop", "", "")
}

func (l *ChatLogic) writeSSE(w http.ResponseWriter, flusher http.Flusher,
	eventType, content, toolName, toolData string) {
	data := map[string]interface{}{
		"type":    eventType,
		"content": content,
	}
	if toolName != "" {
		data["tool_name"] = toolName
		data["tool_data"] = toolData
	}
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, string(jsonData))
	flusher.Flush()
}
