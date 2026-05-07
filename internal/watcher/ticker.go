// Package watcher provides the core loop that periodically checks
// whether monitored cron jobs have run within their expected schedule.
package watcher

import (
	"context"
	"time"
)

// Ticker abstracts time.Ticker so the watcher loop can be tested
// without real wall-clock delays.
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// realTicker wraps a stdlib *time.Ticker.
type realTicker struct {
	t *time.Ticker
}

// NewRealTicker returns a Ticker backed by the standard library.
func NewRealTicker(d time.Duration) Ticker {
	return &realTicker{t: time.NewTicker(d)}
}

func (r *realTicker) C() <-chan time.Time { return r.t.C }
func (r *realTicker) Stop()              { r.t.Stop() }

// fakeTicker is a controllable Ticker used in tests.
type fakeTicker struct {
	ch   chan time.Time
	done chan struct{}
}

// newFakeTicker creates a fakeTicker whose channel can be driven manually.
func newFakeTicker() *fakeTicker {
	return &fakeTicker{
		ch:   make(chan time.Time, 1),
		done: make(chan struct{}),
	}
}

func (f *fakeTicker) C() <-chan time.Time { return f.ch }
func (f *fakeTicker) Stop()              { close(f.done) }

// tick sends a single tick with the given time, respecting context cancellation.
func (f *fakeTicker) tick(ctx context.Context, at time.Time) {
	select {
	case f.ch <- at:
	case <-ctx.Done():
	}
}
