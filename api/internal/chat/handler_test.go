package chat_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func postChat(t *testing.T, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	handler := chat.NewHandler(&mockProvider{response: "ok"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func marshalMessages(t *testing.T, msgs []chat.Message) []byte {
	t.Helper()
	body, err := json.Marshal(chat.ChatRequest{Messages: msgs})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return body
}

func TestChatHandler_RejectsOversizedBody(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":"` + strings.Repeat("a", chat.MaxBodyBytes) + `"}]}`)
	rec := postChat(t, body)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestChatHandler_RejectsInvalidRole(t *testing.T) {
	for _, role := range []string{"system", "model", "", "User"} {
		rec := postChat(t, marshalMessages(t, []chat.Message{{Role: role, Content: "hi"}}))
		if rec.Code != http.StatusBadRequest {
			t.Errorf("role %q: expected 400, got %d", role, rec.Code)
		}
	}
}

func TestChatHandler_RejectsLongUserMessage(t *testing.T) {
	msg := chat.Message{Role: "user", Content: strings.Repeat("a", chat.MaxUserContentRunes+1)}
	rec := postChat(t, marshalMessages(t, []chat.Message{msg}))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestChatHandler_RejectsLongAssistantMessage(t *testing.T) {
	msgs := []chat.Message{
		{Role: "user", Content: "hi"},
		{Role: "assistant", Content: strings.Repeat("a", chat.MaxAssistantContentRunes+1)},
		{Role: "user", Content: "hi"},
	}
	rec := postChat(t, marshalMessages(t, msgs))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// The largest conversation the ChatWidget can send must pass every limit,
// even with 4-byte characters and JSON escaping.
func TestChatHandler_AcceptsLargestWidgetConversation(t *testing.T) {
	msgs := make([]chat.Message, 20)
	for i := range msgs {
		if i%2 == 0 {
			msgs[i] = chat.Message{Role: "user", Content: strings.Repeat("😀", chat.MaxUserContentRunes)}
		} else {
			msgs[i] = chat.Message{Role: "assistant", Content: strings.Repeat("\x01", chat.MaxAssistantContentRunes)}
		}
	}
	body := marshalMessages(t, msgs)
	if len(body) > chat.MaxBodyBytes {
		t.Fatalf("worst-case widget body is %d bytes, over the %d byte cap", len(body), chat.MaxBodyBytes)
	}
	rec := postChat(t, body)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestChatHandler_LimitsCountRunesNotBytes(t *testing.T) {
	// MaxUserContentRunes multi-byte characters is within the limit.
	msg := chat.Message{Role: "user", Content: strings.Repeat("é", chat.MaxUserContentRunes)}
	rec := postChat(t, marshalMessages(t, []chat.Message{msg}))
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}
