package chat_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

func TestChatHandler_ReturnsResponse(t *testing.T) {
	provider := &mockProvider{response: "Rafael is a Software Engineer based in Toronto."}
	handler := chat.NewHandler(provider)

	reqBody := chat.ChatRequest{
		Messages: []chat.Message{{Role: "user", Content: "Who is Rafael?"}},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp chat.ChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Message == "" {
		t.Error("expected non-empty message in response")
	}
}

func TestChatHandler_RejectsEmptyMessages(t *testing.T) {
	provider := &mockProvider{response: "test"}
	handler := chat.NewHandler(provider)

	reqBody := chat.ChatRequest{Messages: []chat.Message{}}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestChatHandler_EnforcesMessageLimit(t *testing.T) {
	provider := &mockProvider{response: "test"}
	handler := chat.NewHandler(provider)

	msgs := make([]chat.Message, 21)
	for i := range msgs {
		msgs[i] = chat.Message{Role: "user", Content: "hello"}
	}

	reqBody := chat.ChatRequest{Messages: msgs}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for too many messages, got %d", rec.Code)
	}
}
