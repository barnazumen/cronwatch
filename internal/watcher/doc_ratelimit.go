// Package watcher provides the core polling loop and supporting utilities
// for cronwatch.
//
// # RateLimit
//
// RateLimit prevents alert storms by suppressing repeated alerts for the same
// job within a configurable cooldown window.
//
// Usage:
//
//	rl := watcher.NewRateLimit(5 * time.Minute)
//
//	if rl.Allow(job.Name()) {
//		// forward the alert
//	}
//
// When a job recovers, call Reset to clear the suppression state so the next
// failure is reported immediately:
//
//	rl.Reset(job.Name())
//
// RateLimit is safe for concurrent use.
package watcher
