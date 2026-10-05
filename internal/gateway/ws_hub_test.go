package gateway

import (
	"testing"
)

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := NewHub()

	client := &Client{
		UserID: 1001,
		Send:   make(chan []byte, 10),
	}

	hub.Register(client)

	if !hub.IsOnline(1001) {
		t.Fatalf("expected user 1001 online")
	}

	hub.BroadcastToUser(1001, []byte("hello"))

	select {
	case msg := <-client.Send:
		if string(msg) != "hello" {
			t.Fatalf("expected hello, got %s", string(msg))
		}
	default:
		t.Fatalf("no message received")
	}

	hub.Unregister(client)
	if hub.IsOnline(1001) {
		t.Fatalf("expected user 1001 offline after unregister")
	}
}
