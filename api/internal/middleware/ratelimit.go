package middleware

import (
	"container/list"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

const (
	defaultWindow  = time.Minute
	defaultMaxKeys = 10000
)

// RateLimitConfig configures RateLimit. Only Limit is required.
type RateLimitConfig struct {
	// Limit is the number of requests one client may make per window. Must be > 0.
	Limit int
	// Window is the length of each fixed window. Defaults to one minute.
	Window time.Duration
	// MaxKeys bounds how many clients are tracked at once. Defaults to 10,000.
	MaxKeys int
	// TrustFlyClientIP keys clients on the Fly-Client-IP header. Enable it only
	// when every request reaches the app through Fly Proxy.
	TrustFlyClientIP bool
	// Now is the clock; defaults to time.Now. It exists for tests.
	Now func() time.Time
}

type entry struct {
	key   string
	count int
	start time.Time
}

type limiter struct {
	limit   int
	window  time.Duration
	maxKeys int
	now     func() time.Time

	mu      sync.Mutex
	entries map[string]*list.Element
	// order holds entries by window start, oldest at the front. A window
	// length is fixed, so the front is always the next to expire.
	order *list.List
}

// RateLimit allows each client Limit requests per fixed window.
//
// A client's window starts at its first request and ends Window later, with
// the count reset then. Rejected requests change nothing, so a client that
// keeps retrying cannot extend its own block. A client can send up to 2×Limit
// requests across a window boundary; that is acceptable for abuse protection.
//
// At most MaxKeys clients are tracked. Expired windows are dropped first. If
// the map is still full, the client with the oldest window is dropped, which
// resets that client's count. This fails open for that one client rather than
// closed for everyone: rejecting new clients when full would let an attacker
// with MaxKeys addresses (a single IPv6 /48 holds 65,536 /64 prefixes) lock
// every new visitor out indefinitely by renewing its entries. Memory stays
// bounded either way; the daily Gemini cap and Fly's concurrency limit are the
// backstops against a flood.
func RateLimit(cfg RateLimitConfig) func(http.Handler) http.Handler {
	if cfg.Limit <= 0 {
		panic(fmt.Sprintf("middleware: RateLimit limit must be positive, got %d", cfg.Limit))
	}
	l := &limiter{
		limit:   cfg.Limit,
		window:  cfg.Window,
		maxKeys: cfg.MaxKeys,
		now:     cfg.Now,
		entries: make(map[string]*list.Element),
		order:   list.New(),
	}
	if l.window <= 0 {
		l.window = defaultWindow
	}
	if l.maxKeys <= 0 {
		l.maxKeys = defaultMaxKeys
	}
	if l.now == nil {
		l.now = time.Now
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ok, retryAfter := l.allow(clientKey(r, cfg.TrustFlyClientIP)); !ok {
				WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded", retryAfter)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// allow records a request for key. When it refuses, it returns how long until
// a retry can succeed.
func (l *limiter) allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()

	for front := l.order.Front(); front != nil; front = l.order.Front() {
		e := front.Value.(*entry)
		if now.Sub(e.start) < l.window {
			break
		}
		l.remove(front)
	}

	if el, tracked := l.entries[key]; tracked {
		e := el.Value.(*entry)
		if e.count >= l.limit {
			return false, e.start.Add(l.window).Sub(now)
		}
		e.count++
		return true, 0
	}

	if len(l.entries) >= l.maxKeys {
		l.remove(l.order.Front())
	}
	l.entries[key] = l.order.PushBack(&entry{key: key, count: 1, start: now})
	return true, 0
}

func (l *limiter) remove(el *list.Element) {
	delete(l.entries, el.Value.(*entry).key)
	l.order.Remove(el)
}

// WriteJSONError writes the API's {"error": "..."} body with a JSON content
// type and a Retry-After header rounded up to whole seconds (minimum 1).
func WriteJSONError(w http.ResponseWriter, status int, msg string, retryAfter time.Duration) {
	secs := int(math.Ceil(retryAfter.Seconds()))
	if secs < 1 {
		secs = 1
	}
	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Retry-After", fmt.Sprint(secs))
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
