package gateway

import (
	"sync"
)

// Client 表示单个 WebSocket 客户端连接
type Client struct {
	UserID int64
	Send   chan []byte
}

// Hub 维护单机在线长连接
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]*Client
}

// NewHub 创建连接管理器
func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]*Client),
	}
}

// Register 客户端上线注册
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.UserID] = c
}

// Unregister 客户端下线注销
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if current, exists := h.clients[c.UserID]; exists && current == c {
		delete(h.clients, c.UserID)
		close(c.Send)
	}
}

// IsOnline 检查用户是否本地在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.clients[userID]
	return exists
}

// BroadcastToUser 向指定用户下发消息
func (h *Hub) BroadcastToUser(userID int64, msg []byte) bool {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		return false
	}

	select {
	case client.Send <- msg:
		return true
	default:
		return false
	}
}

// DefaultHub 全局单例连接管理器
var DefaultHub = NewHub()
