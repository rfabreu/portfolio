package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rfabreu/portfolio-api/internal/middleware"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func newFakeClock() *fakeClock { return &fakeClock{t: time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)} }

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

var ok200 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

func limited(cfg middleware.RateLimitConfig) http.Handler {
	return middleware.RateLimit(cfg)(ok200)
}

func do(h http.Handler, remote string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", nil)
	req.RemoteAddr = remote
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRateLimiter_AllowsUpToLimitThenBlocks(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 3, Now: clk.Now})
	for i := 0; i < 3; i++ {
		if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
	}
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusTooManyRequests {
		t.Errorf("request over the limit: expected 429, got %d", rec.Code)
	}
	if rec := do(h, "192.168.1.2:1", nil); rec.Code != http.StatusOK {
		t.Errorf("other client must have its own bucket, got %d", rec.Code)
	}
}

func TestRateLimiter_WindowResetsOnItsOwn(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "192.168.1.1:1", nil)
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	clk.Advance(time.Minute)
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusOK {
		t.Errorf("window must reset after 60s, got %d", rec.Code)
	}
}

// The old limiter refreshed lastSeen on rejected requests, so a client
// retrying every 59s stayed blocked forever.
func TestRateLimiter_RejectedRequestsDoNotExtendWindow(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "192.168.1.1:1", nil) // window starts at t0
	for i := 0; i < 5; i++ {
		clk.Advance(10 * time.Second)
		if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusTooManyRequests {
			t.Fatalf("t0+%ds: expected 429, got %d", (i+1)*10, rec.Code)
		}
	}
	clk.Advance(10 * time.Second) // t0+60s
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusOK {
		t.Errorf("window must end at t0+60s regardless of rejected traffic, got %d", rec.Code)
	}
}

func TestRateLimiter_WindowBoundaryIsExact(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "192.168.1.1:1", nil)
	clk.Advance(time.Minute - time.Nanosecond)
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusTooManyRequests {
		t.Errorf("1ns before the window ends: expected 429, got %d", rec.Code)
	}
	clk.Advance(time.Nanosecond)
	if rec := do(h, "192.168.1.1:1", nil); rec.Code != http.StatusOK {
		t.Errorf("exactly at the window end: expected 200, got %d", rec.Code)
	}
}

func TestRateLimiter_429IsJSONWithRetryAfter(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "192.168.1.1:1", nil)
	clk.Advance(20 * time.Second)
	rec := do(h, "192.168.1.1:1", nil)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if got := rec.Header().Get("Retry-After"); got != "40" {
		t.Errorf("Retry-After = %q, want 40", got)
	}
	if body := rec.Body.String(); body != `{"error":"rate limit exceeded"}`+"\n" {
		t.Errorf("body = %q", body)
	}
}

func TestRateLimiter_RetryAfterRoundsUpAndIsAtLeastOne(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "192.168.1.1:1", nil)
	clk.Advance(time.Minute - 100*time.Millisecond)
	rec := do(h, "192.168.1.1:1", nil)
	if got := rec.Header().Get("Retry-After"); got != "1" {
		t.Errorf("Retry-After = %q, want 1", got)
	}
}

func TestRateLimiter_FlyClientIPSeparatesBucketsWhenTrusted(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, TrustFlyClientIP: true, Now: clk.Now})
	remote := "172.16.0.1:1000"
	if rec := do(h, remote, map[string]string{"Fly-Client-IP": "203.0.113.1"}); rec.Code != 200 {
		t.Fatalf("got %d", rec.Code)
	}
	if rec := do(h, remote, map[string]string{"Fly-Client-IP": "203.0.113.2"}); rec.Code != 200 {
		t.Errorf("different Fly-Client-IP must get its own bucket, got %d", rec.Code)
	}
	if rec := do(h, remote, map[string]string{"Fly-Client-IP": "203.0.113.1"}); rec.Code != 429 {
		t.Errorf("same Fly-Client-IP must share a bucket, got %d", rec.Code)
	}
}

func TestRateLimiter_ForgedFlyClientIPIgnoredWhenNotTrusted(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 2, Now: clk.Now})
	var codes []int
	for i := 0; i < 4; i++ {
		rec := do(h, "172.16.0.1:1000", map[string]string{"Fly-Client-IP": fmt.Sprintf("203.0.113.%d", i)})
		codes = append(codes, rec.Code)
	}
	want := []int{200, 200, 429, 429}
	for i := range want {
		if codes[i] != want[i] {
			t.Fatalf("codes = %v, want %v", codes, want)
		}
	}
}

func TestRateLimiter_XForwardedForNeverUsed(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, TrustFlyClientIP: true, Now: clk.Now})
	do(h, "172.16.0.1:1000", map[string]string{"X-Forwarded-For": "203.0.113.1"})
	if rec := do(h, "172.16.0.1:1000", map[string]string{"X-Forwarded-For": "203.0.113.2"}); rec.Code != 429 {
		t.Errorf("rotating X-Forwarded-For must not bypass the limit, got %d", rec.Code)
	}
}

// When the map is full the oldest window is dropped, so a flood of new keys
// can neither grow memory nor lock new visitors out.
func TestRateLimiter_FullMapEvictsOldestWindow(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, MaxKeys: 2, Now: clk.Now})

	do(h, "10.0.0.1:1", nil) // oldest
	clk.Advance(time.Second)
	do(h, "10.0.0.2:1", nil)
	if rec := do(h, "10.0.0.1:1", nil); rec.Code != 429 {
		t.Fatalf("10.0.0.1 is over its limit, got %d", rec.Code)
	}
	clk.Advance(time.Second)
	if rec := do(h, "10.0.0.3:1", nil); rec.Code != 200 {
		t.Fatalf("a new key must be admitted while the map is full, got %d", rec.Code)
	}
	if rec := do(h, "10.0.0.2:1", nil); rec.Code != 429 {
		t.Errorf("10.0.0.2 was not evicted and is still limited, got %d", rec.Code)
	}
	// 10.0.0.1 had the oldest window, so it was evicted and starts afresh.
	if rec := do(h, "10.0.0.1:1", nil); rec.Code != 200 {
		t.Errorf("evicted key starts a new window, got %d", rec.Code)
	}
}

func TestRateLimiter_NewKeysAreNeverLockedOutDuringAFlood(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, MaxKeys: 3, Now: clk.Now})
	for i := 0; i < 100; i++ {
		clk.Advance(time.Second)
		if rec := do(h, fmt.Sprintf("10.0.%d.1:1", i), nil); rec.Code != 200 {
			t.Fatalf("new key %d: got %d", i, rec.Code)
		}
	}
}

func TestRateLimiter_ExpiredWindowsAreForgotten(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, MaxKeys: 2, Now: clk.Now})
	do(h, "10.0.0.1:1", nil)
	do(h, "10.0.0.2:1", nil)
	clk.Advance(time.Minute)
	// Both windows expired: both keys are fresh, and neither displaces the other.
	for _, ip := range []string{"10.0.0.1:1", "10.0.0.2:1"} {
		if rec := do(h, ip, nil); rec.Code != 200 {
			t.Errorf("%s after expiry: got %d", ip, rec.Code)
		}
	}
	for _, ip := range []string{"10.0.0.1:1", "10.0.0.2:1"} {
		if rec := do(h, ip, nil); rec.Code != 429 {
			t.Errorf("%s second request in the new window: got %d", ip, rec.Code)
		}
	}
}

func TestRateLimiter_DefaultsToOneMinuteWindow(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Now: clk.Now})
	do(h, "10.0.0.1:1", nil)
	clk.Advance(59 * time.Second)
	if rec := do(h, "10.0.0.1:1", nil); rec.Code != 429 {
		t.Errorf("got %d", rec.Code)
	}
}

func TestRateLimiter_CustomWindow(t *testing.T) {
	clk := newFakeClock()
	h := limited(middleware.RateLimitConfig{Limit: 1, Window: time.Hour, Now: clk.Now})
	do(h, "10.0.0.1:1", nil)
	clk.Advance(30 * time.Minute)
	rec := do(h, "10.0.0.1:1", nil)
	if rec.Code != 429 || rec.Header().Get("Retry-After") != strconv.Itoa(1800) {
		t.Errorf("got %d Retry-After %q", rec.Code, rec.Header().Get("Retry-After"))
	}
}

func TestRateLimiter_RejectsNonPositiveLimit(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for Limit <= 0")
		}
	}()
	middleware.RateLimit(middleware.RateLimitConfig{Limit: 0})
}

func TestRateLimiter_ConcurrentUse(t *testing.T) {
	h := limited(middleware.RateLimitConfig{Limit: 50})
	var wg sync.WaitGroup
	var okCount atomic.Int32
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if do(h, "10.0.0.1:1", nil).Code == 200 {
				okCount.Add(1)
			}
		}()
	}
	wg.Wait()
	if okCount.Load() != 50 {
		t.Errorf("served %d, want exactly 50", okCount.Load())
	}
}
