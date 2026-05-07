package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/watcher"
)

// TestHistory_ConcurrentWrites verifies that concurrent Record calls do not
// race or corrupt internal state.
func TestHistory_ConcurrentWrites(t *testing.T) {
	h := watcher.NewHistory(50)
	done := make(chan struct{})

	write := func() {
		for i := 0; i < 30; i++ {
			h.Record(watcher.CheckRecord{
				At:      time.Now(),
				Missed:  []string{"job"},
				Alerted: []string{"job"},
			})
		}
		done <- struct{}{}
	}

	const goroutines = 5
	for i := 0; i < goroutines; i++ {
		go write()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}

	if h.Len() > 50 {
		t.Errorf("history exceeded cap: len=%d", h.Len())
	}
}

// TestHistory_SnapshotWhileWriting checks that Snapshot can be called safely
// while another goroutine is recording.
func TestHistory_SnapshotWhileWriting(t *testing.T) {
	h := watcher.NewHistory(20)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				h.Record(watcher.CheckRecord{At: time.Now()})
			}
		}
	}()

	for i := 0; i < 20; i++ {
		_ = h.Snapshot()
		time.Sleep(time.Millisecond)
	}
	close(stop)
}
