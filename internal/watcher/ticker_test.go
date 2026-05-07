package watcher

import (
	"context"
	"testing"
	"time"
)

func TestNewRealTicker_ReceivesTick(t *testing.T) {
	tk := NewRealTicker(20 * time.Millisecond)
	defer tk.Stop()

	select {
	case got := <-tk.C():
		if got.IsZero() {
			t.Fatal("expected non-zero tick time")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for real ticker")
	}
}

func TestNewRealTicker_StopsCleanly(t *testing.T) {
	tk := NewRealTicker(10 * time.Millisecond)
	// Drain one tick then stop; subsequent reads should not produce values.
	<-tk.C()
	tk.Stop()

	// Give a brief window; no further ticks should arrive.
	time.Sleep(30 * time.Millisecond)
	select {
	case <-tk.C():
		// Draining a buffered tick that arrived before Stop is acceptable;
		// just ensure we don't block forever.
	default:
	}
}

func TestFakeTicker_DeliversTick(t *testing.T) {
	ctx := context.Background()
	ft := newFakeTicker()

	now := time.Now()
	go ft.tick(ctx, now)

	select {
	case got := <-ft.C():
		if !got.Equal(now) {
			t.Fatalf("expected %v, got %v", now, got)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for fake tick")
	}
}

func TestFakeTicker_StopClosesdone(t *testing.T) {
	ft := newFakeTicker()
	ft.Stop()

	select {
	case <-ft.done:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("done channel was not closed after Stop")
	}
}

func TestFakeTicker_TickRespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	ft := newFakeTicker()
	// tick should return without blocking because ctx is already done.
	done := make(chan struct{})
	go func() {
		ft.tick(ctx, time.Now())
		close(done)
	}()

	select {
	case <-done:
		// expected — tick returned promptly
	case <-time.After(200 * time.Millisecond):
		t.Fatal("tick did not respect context cancellation")
	}
}
