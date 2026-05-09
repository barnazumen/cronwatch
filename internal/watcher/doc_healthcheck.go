// Package watcher provides the core monitoring loop for cronwatch.
//
// # HealthCheck
//
// HealthCheck maintains the current health state for each monitored job.
// It is safe for concurrent use.
//
// Usage:
//
//	hc := watcher.NewHealthCheck()
//
//	// After a successful run:
//	hc.RecordHealthy("nightly-backup")
//
//	// After a missed or failed run:
//	hc.RecordUnhealthy("nightly-backup", "missed scheduled run")
//
//	// Query a single job:
//	if s, ok := hc.Status("nightly-backup"); ok && !s.Healthy {
//		// handle degraded state
//	}
//
//	// Dump all statuses (e.g. for an HTTP /healthz endpoint):
//	all := hc.Snapshot()
package watcher
