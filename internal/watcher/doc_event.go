// Package watcher — event.go
//
// # Event
//
// An Event represents a single observation made by the watcher during one
// check cycle. Every event has a Kind (missed, failure, or recovered), a
// reference to the affected Job, and a timestamp.
//
// Events are produced inside the watcher loop and forwarded to the alerter
// pipeline. Downstream consumers (log alerter, email alerter, etc.) receive
// an Event and decide how to surface it to operators.
//
// # Usage
//
//	ev := watcher.NewEvent(watcher.EventMissed, j)
//	fmt.Println(ev.Message)
//
// # Event kinds
//
//   - EventMissed   — job did not run within its grace window
//   - EventFailure  — job ran but exited non-zero
//   - EventRecovered — job is healthy again after a prior alert
package watcher
