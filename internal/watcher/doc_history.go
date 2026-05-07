// Package watcher provides the core loop that checks cron jobs for missed runs.
//
// # History
//
// History maintains a bounded, thread-safe ring-buffer of [CheckRecord] values
// produced by each watcher tick cycle.
//
// Each [CheckRecord] captures:
//   - At:      the wall-clock time the check ran
//   - Missed:  names of jobs whose last-seen time exceeded their schedule window
//   - Alerted: names of jobs for which an alert was successfully dispatched
//
// Usage:
//
//	h := watcher.NewHistory(200)
//
//	// inside the watcher tick:
//	h.Record(watcher.CheckRecord{
//		At:      time.Now(),
//		Missed:  missedNames,
//		Alerted: alertedNames,
//	})
//
//	// inspect recent activity:
//	for _, rec := range h.Snapshot() {
//		fmt.Println(rec.At, rec.Missed)
//	}
package watcher
