package chat_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []chat.Message) (string, error) {
	return m.response, m.err
}

func TestMockProvider_ReturnsResponse(t *testing.T) {
	provider := &mockProvider{response: "Hello! I'm Rafael's AI assistant."}
	resp, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{
		{Role: "user", Content: "Hi"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Hello! I'm Rafael's AI assistant." {
		t.Errorf("expected greeting, got: %s", resp)
	}
}

const testAPIKey = "test-secret-key-do-not-log"

func TestGeminiProvider_ConnectionErrorDoesNotLeakKey(t *testing.T) {
	// Start and immediately close a server so the address refuses connections.
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()

	provider := chat.NewGeminiProvider(testAPIKey, "")
	provider.BaseURL = srv.URL

	_, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected an error for a refused connection")
	}
	if strings.Contains(err.Error(), testAPIKey) {
		t.Errorf("error leaks the API key: %v", err)
	}
}

func TestGeminiProvider_APIErrorDoesNotLeakKey(t *testing.T) {
	var gotHeader, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("x-goog-api-key")
		gotQuery = r.URL.RawQuery
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
	}))
	defer srv.Close()

	provider := chat.NewGeminiProvider(testAPIKey, "")
	provider.BaseURL = srv.URL

	_, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected an error for a 500 response")
	}
	if strings.Contains(err.Error(), testAPIKey) {
		t.Errorf("error leaks the API key: %v", err)
	}
	if gotHeader != testAPIKey {
		t.Errorf("expected key in x-goog-api-key header, got %q", gotHeader)
	}
	if strings.Contains(gotQuery, testAPIKey) {
		t.Errorf("API key must not be sent in the URL query, got %q", gotQuery)
	}
}

func TestGeminiProvider_ReturnsText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models/gemini-2.5-flash-lite:generateContent" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"candidates":[{"content":{"parts":[{"text":"Hello"}]}}]}`)
	}))
	defer srv.Close()

	provider := chat.NewGeminiProvider(testAPIKey, "")
	provider.BaseURL = srv.URL

	resp, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{{Role: "user", Content: "Hi"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Hello" {
		t.Errorf("expected Hello, got %q", resp)
	}
}

func TestGeminiProvider_TimesOut(t *testing.T) {
	if chat.NewGeminiProvider(testAPIKey, "").HTTPClient.Timeout <= 0 {
		t.Fatal("default HTTP client must have a timeout")
	}

	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	provider := chat.NewGeminiProvider(testAPIKey, "")
	provider.BaseURL = srv.URL
	provider.HTTPClient = &http.Client{Timeout: 50 * time.Millisecond}

	_, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if strings.Contains(err.Error(), testAPIKey) {
		t.Errorf("error leaks the API key: %v", err)
	}
}
