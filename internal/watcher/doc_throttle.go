// Package watcher — Throttle
//
// Throttle enforces a minimum interval between consecutive alert firings for a
// given key (typically a job name). It is a simpler, stateless alternative to
// RateLimit when you only need to prevent alert storms without tracking
// cooldown windows explicitly.
//
// Usage:
//
//	th := watcher.NewThrottle(5 * time.Minute)
//
//	if th.Allow(job.Name) {
//		// fire the alert
//	}
//
// Throttle is safe for concurrent use by multiple goroutines.
package watcher
