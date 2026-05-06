// Package watcher provides a background loop that polls the job registry
// and emits notifications when a cron job has exceeded its maximum allowed
// interval without recording a successful run.
//
// Usage:
//
//	notify := make(chan *job.Job, 16)
//	w := watcher.New(registry, 30*time.Second, notify)
//	w.Start()
//	defer w.Stop()
//
//	for missed := range notify {
//		// handle missed job alert
//	}
package watcher
