package middleware

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

func TestLimiter_MapAndQueueNeverExceedMaxKeys(t *testing.T) {
	now := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	l := &limiter{
		limit: 2, window: time.Minute, maxKeys: 50,
		now:     func() time.Time { return now },
		entries: map[string]*list.Element{},
		order:   list.New(),
	}
	for i := 0; i < 5000; i++ {
		now = now.Add(100 * time.Millisecond)
		l.allow(fmt.Sprintf("key-%d", i%700))
		if len(l.entries) > l.maxKeys || l.order.Len() != len(l.entries) {
			t.Fatalf("step %d: %d map entries, %d queued, max %d", i, len(l.entries), l.order.Len(), l.maxKeys)
		}
	}
}

func TestLimiter_QueueStaysOrderedByStart(t *testing.T) {
	now := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	l := &limiter{
		limit: 1, window: time.Minute, maxKeys: 10,
		now:     func() time.Time { return now },
		entries: map[string]*list.Element{},
		order:   list.New(),
	}
	for i := 0; i < 500; i++ {
		now = now.Add(7 * time.Second)
		l.allow(fmt.Sprintf("key-%d", i%13))
	}
	var prev time.Time
	for el := l.order.Front(); el != nil; el = el.Next() {
		e := el.Value.(*entry)
		if e.start.Before(prev) {
			t.Fatalf("queue out of order: %v before %v", e.start, prev)
		}
		prev = e.start
	}
}
