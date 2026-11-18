package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

// fakeProvider replaces Gemini; tests never reach the network.
type fakeProvider struct{ calls atomic.Int32 }

func (p *fakeProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []chat.Message) (string, error) {
	p.calls.Add(1)
	return "ok", nil
}

const (
	chatBody = `{"messages":[{"role":"user","content":"hi"}]}`
	gameBody = `{"board":["X","","","","","","","",""],"session":{"wins":0,"losses":0,"draws":0}}`
)

func testConfig() config {
	return config{allowedOrigins: "http://localhost:4321", chatRPM: 2, gameRPM: 3, chatClientDaily: 100, chatDailyCap: 100}
}

func send(h http.Handler, method, path, body, flyIP string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.RemoteAddr = "172.16.0.1:1234"
	if flyIP != "" {
		req.Header.Set("Fly-Client-IP", flyIP)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestNewHandler_HealthIsNeverLimited(t *testing.T) {
	h := newHandler(testConfig(), &fakeProvider{})
	for i := 0; i < 200; i++ {
		if rec := send(h, http.MethodGet, "/health", "", ""); rec.Code != http.StatusOK {
			t.Fatalf("health request %d: %d", i, rec.Code)
		}
	}
}

func TestNewHandler_RoutesHaveIndependentLimits(t *testing.T) {
	h := newHandler(testConfig(), &fakeProvider{})
	for i := 0; i < 2; i++ {
		if rec := send(h, http.MethodPost, "/api/chat", chatBody, ""); rec.Code != http.StatusOK {
			t.Fatalf("chat %d: %d", i, rec.Code)
		}
	}
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, ""); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("chat over its limit: %d", rec.Code)
	}
	// Chat being exhausted must not touch game, and /health stays open.
	for i := 0; i < 3; i++ {
		if rec := send(h, http.MethodPost, "/api/game", gameBody, ""); rec.Code != http.StatusOK {
			t.Fatalf("game %d: %d %s", i, rec.Code, rec.Body.String())
		}
	}
	if rec := send(h, http.MethodPost, "/api/game", gameBody, ""); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("game over its limit: %d", rec.Code)
	}
	if rec := send(h, http.MethodGet, "/health", "", ""); rec.Code != http.StatusOK {
		t.Fatalf("health: %d", rec.Code)
	}
}

func TestNewHandler_LimitResponseIsJSONWithRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	cfg := testConfig()
	cfg.now = func() time.Time { return now }
	h := newHandler(cfg, &fakeProvider{})
	var rec *httptest.ResponseRecorder
	for i := 0; i < 3; i++ {
		rec = send(h, http.MethodPost, "/api/chat", chatBody, "")
	}
	if rec.Code != http.StatusTooManyRequests ||
		rec.Header().Get("Content-Type") != "application/json" ||
		rec.Header().Get("Retry-After") != "60" {
		t.Errorf("got %d %v", rec.Code, rec.Header())
	}
}

func TestNewHandler_FlyClientIPOnlyHonouredWhenTrusted(t *testing.T) {
	cfg := testConfig()
	cfg.chatRPM = 1

	cfg.trustFlyClientIP = true
	h := newHandler(cfg, &fakeProvider{})
	send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.1")
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.2"); rec.Code != http.StatusOK {
		t.Errorf("trusted: a different client must get its own bucket, got %d", rec.Code)
	}

	cfg.trustFlyClientIP = false
	h = newHandler(cfg, &fakeProvider{})
	send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.1")
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.2"); rec.Code != http.StatusTooManyRequests {
		t.Errorf("untrusted: a forged header must not get a fresh bucket, got %d", rec.Code)
	}
}

func TestNewHandler_DailyCapReturns503AndStopsCallingProvider(t *testing.T) {
	cfg := testConfig()
	cfg.chatRPM = 100
	cfg.chatDailyCap = 3
	cfg.trustFlyClientIP = true
	provider := &fakeProvider{}
	h := newHandler(cfg, provider)

	// Different clients, so only the global cap can stop them.
	for i, ip := range []string{"203.0.113.1", "203.0.113.2", "203.0.113.3"} {
		if rec := send(h, http.MethodPost, "/api/chat", chatBody, ip); rec.Code != http.StatusOK {
			t.Fatalf("call %d: %d", i, rec.Code)
		}
	}
	rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.4")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" || rec.Header().Get("Retry-After") == "" {
		t.Errorf("headers = %v", rec.Header())
	}
	if provider.calls.Load() != 3 {
		t.Errorf("provider called %d times, want 3", provider.calls.Load())
	}
	// The game is unaffected by the chat cap.
	if rec := send(h, http.MethodPost, "/api/game", gameBody, "203.0.113.4"); rec.Code != http.StatusOK {
		t.Errorf("game while chat cap is spent: %d", rec.Code)
	}
}

func TestNewHandler_PerClientLimitRunsBeforeProvider(t *testing.T) {
	provider := &fakeProvider{}
	h := newHandler(testConfig(), provider) // chat limit 2
	for i := 0; i < 10; i++ {
		send(h, http.MethodPost, "/api/chat", chatBody, "")
	}
	if provider.calls.Load() != 2 {
		t.Errorf("provider called %d times, want 2: limited requests must not reach it", provider.calls.Load())
	}
}

func TestNewHandler_PreflightBypassesLimiter(t *testing.T) {
	h := newHandler(testConfig(), &fakeProvider{})
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodOptions, "/api/chat", nil)
		req.Header.Set("Origin", "http://localhost:4321")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("preflight %d: %d", i, rec.Code)
		}
	}
}

func TestNewHandler_LimitResponsesCarryCORSHeaders(t *testing.T) {
	// The widget reads the status of a 429, so the browser must be allowed to see it.
	h := newHandler(testConfig(), &fakeProvider{})
	var rec *httptest.ResponseRecorder
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader([]byte(chatBody)))
		req.Header.Set("Origin", "http://localhost:4321")
		req.RemoteAddr = "172.16.0.1:1234"
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:4321" {
		t.Errorf("got %d %v", rec.Code, rec.Header())
	}
}

func TestNewHandler_UsesInjectedClock(t *testing.T) {
	now := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	cfg := testConfig()
	cfg.chatRPM = 1
	cfg.now = func() time.Time { return now }
	h := newHandler(cfg, &fakeProvider{})
	send(h, http.MethodPost, "/api/chat", chatBody, "")
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, ""); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("got %d", rec.Code)
	}
	now = now.Add(time.Minute)
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, ""); rec.Code != http.StatusOK {
		t.Errorf("after the window: %d", rec.Code)
	}
}

func TestNewHandler_PerClientDailyLimitStopsOneClientDrainingTheCap(t *testing.T) {
	now := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	cfg := testConfig()
	cfg.chatRPM = 100
	cfg.chatClientDaily = 3
	cfg.chatDailyCap = 100
	cfg.trustFlyClientIP = true
	cfg.now = func() time.Time { return now }
	provider := &fakeProvider{}
	h := newHandler(cfg, provider)

	for i := 0; i < 3; i++ {
		if rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.1"); rec.Code != http.StatusOK {
			t.Fatalf("call %d: %d", i, rec.Code)
		}
	}
	rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.1")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") != "86400" {
		t.Fatalf("over the daily limit: %d, Retry-After %q", rec.Code, rec.Header().Get("Retry-After"))
	}
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.2"); rec.Code != http.StatusOK {
		t.Errorf("another client is unaffected, got %d", rec.Code)
	}
	if provider.calls.Load() != 4 {
		t.Errorf("provider called %d times, want 4", provider.calls.Load())
	}
	now = now.Add(24 * time.Hour)
	if rec := send(h, http.MethodPost, "/api/chat", chatBody, "203.0.113.1"); rec.Code != http.StatusOK {
		t.Errorf("after 24h the client's budget resets, got %d", rec.Code)
	}
}

func TestNewHandler_CapResponseCarriesCORSHeaders(t *testing.T) {
	cfg := testConfig()
	cfg.chatDailyCap = 1
	h := newHandler(cfg, &fakeProvider{})
	var rec *httptest.ResponseRecorder
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader([]byte(chatBody)))
		req.Header.Set("Origin", "http://localhost:4321")
		req.RemoteAddr = "172.16.0.1:1234"
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
	if rec.Code != http.StatusServiceUnavailable || rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:4321" {
		t.Errorf("got %d %v", rec.Code, rec.Header())
	}
}
