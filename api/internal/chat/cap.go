package chat

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// DailyCapWindow is the length of the window CHAT_DAILY_CAP applies to.
const DailyCapWindow = 24 * time.Hour

// CapExceededError is returned when the Gemini call budget for the current
// window is spent. RetryAfter is the time until the window resets.
type CapExceededError struct {
	RetryAfter time.Duration
}

func (e *CapExceededError) Error() string {
	return fmt.Sprintf("chat call cap reached; resets in %s", e.RetryAfter.Round(time.Second))
}

// CappedProvider limits how many calls reach the wrapped Provider in total,
// across all clients. The count is in memory and per process: it resets when
// the machine restarts, and each machine has its own.
type CappedProvider struct {
	inner  Provider
	limit  int
	window time.Duration
	now    func() time.Time

	mu      sync.Mutex
	count   int
	start   time.Time
	tripped bool
}

// NewCappedProvider allows limit calls per fixed window. The window starts at
// the first call and ends window later, with the count reset then. Rejected
// calls change nothing. A nil now uses time.Now. It panics if limit or window
// is not positive: a zero window would reset on every call and never cap.
func NewCappedProvider(inner Provider, limit int, window time.Duration, now func() time.Time) *CappedProvider {
	if limit <= 0 || window <= 0 {
		panic(fmt.Sprintf("chat: NewCappedProvider needs a positive limit and window, got %d and %s", limit, window))
	}
	if now == nil {
		now = time.Now
	}
	return &CappedProvider{inner: inner, limit: limit, window: window, now: now}
}

func (c *CappedProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	if err := c.take(); err != nil {
		return "", err
	}
	return c.inner.GenerateResponse(ctx, systemPrompt, messages)
}

// take spends one call from the budget. Calls that later fail still count: a
// failed Gemini request may still be billed or counted against quota.
func (c *CappedProvider) take() error {
	now := c.now()
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.count == 0 || now.Sub(c.start) >= c.window {
		c.count, c.start, c.tripped = 0, now, false
	}
	if c.count >= c.limit {
		if !c.tripped {
			c.tripped = true
			log.Printf("WARNING: chat call cap of %d per %s reached; rejecting chat requests until the window resets", c.limit, c.window)
		}
		return &CapExceededError{RetryAfter: c.start.Add(c.window).Sub(now)}
	}
	c.count++
	return nil
}
