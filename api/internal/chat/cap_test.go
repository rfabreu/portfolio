package chat_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

type capClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *capClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *capClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func newCapClock() *capClock { return &capClock{t: time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)} }

// countingProvider stands in for Gemini; tests never reach the network.
type countingProvider struct {
	calls atomic.Int32
	err   error
}

func (p *countingProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []chat.Message) (string, error) {
	p.calls.Add(1)
	if p.err != nil {
		return "", p.err
	}
	return "ok", nil
}

func generate(p chat.Provider) error {
	_, err := p.GenerateResponse(context.Background(), "sys", []chat.Message{{Role: "user", Content: "hi"}})
	return err
}

func TestCappedProvider_StopsCallingInnerAtCap(t *testing.T) {
	inner := &countingProvider{}
	capped := chat.NewCappedProvider(inner, 3, 24*time.Hour, newCapClock().Now)
	for i := 0; i < 3; i++ {
		if err := generate(capped); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	err := generate(capped)
	var capErr *chat.CapExceededError
	if !errors.As(err, &capErr) {
		t.Fatalf("call over the cap: err = %v, want *CapExceededError", err)
	}
	if inner.calls.Load() != 3 {
		t.Errorf("inner provider called %d times, want 3", inner.calls.Load())
	}
}

func TestCappedProvider_WindowResetsAndRejectsDoNotExtendIt(t *testing.T) {
	clk := newCapClock()
	inner := &countingProvider{}
	capped := chat.NewCappedProvider(inner, 1, 24*time.Hour, clk.Now)
	generate(capped) // window starts at t0
	for i := 0; i < 5; i++ {
		clk.Advance(4 * time.Hour)
		if err := generate(capped); err == nil {
			t.Fatalf("t0+%dh: expected cap error", (i+1)*4)
		}
	}
	clk.Advance(4 * time.Hour) // t0+24h
	if err := generate(capped); err != nil {
		t.Errorf("window must reset 24h after it started, got %v", err)
	}
}

func TestCappedProvider_RetryAfterIsTimeToReset(t *testing.T) {
	clk := newCapClock()
	capped := chat.NewCappedProvider(&countingProvider{}, 1, 24*time.Hour, clk.Now)
	generate(capped)
	clk.Advance(6 * time.Hour)
	var capErr *chat.CapExceededError
	if !errors.As(generate(capped), &capErr) {
		t.Fatal("expected cap error")
	}
	if capErr.RetryAfter != 18*time.Hour {
		t.Errorf("RetryAfter = %v, want 18h", capErr.RetryAfter)
	}
}

func TestCappedProvider_FailedCallsStillCount(t *testing.T) {
	inner := &countingProvider{err: errors.New("gemini down")}
	capped := chat.NewCappedProvider(inner, 2, 24*time.Hour, newCapClock().Now)
	generate(capped)
	generate(capped)
	var capErr *chat.CapExceededError
	if !errors.As(generate(capped), &capErr) {
		t.Error("failed Gemini calls must count toward the cap")
	}
}

func TestCappedProvider_ConcurrentCallsNeverExceedCap(t *testing.T) {
	inner := &countingProvider{}
	capped := chat.NewCappedProvider(inner, 25, 24*time.Hour, nil)
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); generate(capped) }()
	}
	wg.Wait()
	if inner.calls.Load() != 25 {
		t.Errorf("inner provider called %d times, want exactly 25", inner.calls.Load())
	}
}

func postChatTo(h http.Handler, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestChatHandler_CapReachedReturns503JSONWithRetryAfter(t *testing.T) {
	clk := newCapClock()
	inner := &countingProvider{}
	h := chat.NewHandler(chat.NewCappedProvider(inner, 1, 24*time.Hour, clk.Now))
	body := marshalMessages(t, []chat.Message{{Role: "user", Content: "hi"}})

	if rec := postChatTo(h, body); rec.Code != http.StatusOK {
		t.Fatalf("first call: %d", rec.Code)
	}
	clk.Advance(time.Hour)
	rec := postChatTo(h, body)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	if ra := rec.Header().Get("Retry-After"); ra != "82800" {
		t.Errorf("Retry-After = %q, want 82800", ra)
	}
	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp["error"] == "" {
		t.Errorf("body = %q, want {\"error\": ...}", rec.Body.String())
	}
}

func TestChatHandler_InvalidRequestsDoNotSpendCap(t *testing.T) {
	inner := &countingProvider{}
	h := chat.NewHandler(chat.NewCappedProvider(inner, 1, 24*time.Hour, newCapClock().Now))
	for i := 0; i < 5; i++ {
		if rec := postChatTo(h, []byte(`{"messages":[]}`)); rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	}
	body := marshalMessages(t, []chat.Message{{Role: "user", Content: "hi"}})
	if rec := postChatTo(h, body); rec.Code != http.StatusOK {
		t.Errorf("a valid request after invalid ones must still pass, got %d", rec.Code)
	}
}

func TestNewCappedProvider_RejectsNonPositiveArguments(t *testing.T) {
	for name, args := range map[string][2]int64{"zero limit": {0, int64(time.Hour)}, "negative limit": {-1, int64(time.Hour)}, "zero window": {1, 0}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: expected a panic", name)
				}
			}()
			chat.NewCappedProvider(&countingProvider{}, int(args[0]), time.Duration(args[1]), nil)
		}()
	}
}
