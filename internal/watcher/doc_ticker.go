// Package watcher — ticker abstraction
//
// # Ticker
//
// The Ticker interface decouples the watcher's check loop from the
// wall clock, making the loop fully testable without sleep-based
// delays.
//
// Production code uses NewRealTicker, which wraps time.Ticker.
// Tests instantiate newFakeTicker and drive ticks manually via
// (*fakeTicker).tick, giving deterministic, fast test execution.
//
// Usage:
//
//	tk := watcher.NewRealTicker(cfg.Interval)
//	defer tk.Stop()
//	for {
//		select {
//		case t := <-tk.C():
//			checkJobs(t)
//		case <-ctx.Done():
//			return
//		}
//	}
package watcher
