package chat_test

import (
	"context"
	"testing"

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
