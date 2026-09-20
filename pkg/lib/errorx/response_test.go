package errorx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpResultSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	HttpResult(req, w, map[string]string{"foo": "bar"}, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp ResponseSuccessBean
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Result != ResultSuccess || resp.Code != OK {
		t.Fatalf("expected success and 200, got %v and %v", resp.Result, resp.Code)
	}
}

func TestHttpResultError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	HttpResult(req, w, nil, NewCodeError(InvalidParam))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp ResponseErrorBean
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Result != ResultFailure || resp.Code != InvalidParam {
		t.Fatalf("expected failure and %d, got %v and %v", InvalidParam, resp.Result, resp.Code)
	}
}
