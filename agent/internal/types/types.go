package types

type ChatRequest struct {
	SessionID string `json:"session_id,optional"`
	Message   string `json:"message"`
}

type HistoryRequest struct {
	SessionID string `form:"session_id"`
}

type HistoryResponse struct {
	SessionID string        `json:"session_id"`
	Title     string        `json:"title"`
	Messages  []MessageItem `json:"messages"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StopRequest struct {
	SessionID string `json:"session_id"`
}

type SessionItem struct {
	SessionID    string `json:"session_id"`
	Title        string `json:"title"`
	MessageCount int32  `json:"message_count"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type ListSessionsResponse struct {
	Sessions []SessionItem `json:"sessions"`
}

type DeleteSessionRequest struct {
	SessionID string `json:"session_id"`
}

type DeleteSessionResponse struct {
	Success bool `json:"success"`
}
