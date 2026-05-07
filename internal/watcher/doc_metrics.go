// Package watcher provides the core polling loop that periodically inspects
// each registered job to detect missed runs and dispatch alerts.
//
// # Metrics
//
// The [Metrics] type exposes atomic counters that accumulate over the lifetime
// of the watcher process:
//
//   - Checks  – total number of per-job evaluations performed
//   - Missed  – jobs whose last-seen timestamp exceeded the configured grace window
//   - Alerted – successful alert dispatches (one per missed detection)
//
// Retrieve a consistent point-in-time view with [Metrics.Snapshot].
//
// # SyslogWriter
//
// [SyslogWriter] emits RFC-3339-timestamped lines compatible with syslog
// forwarders such as rsyslog or Fluentd. Pass nil to [NewSyslogWriter] to
// write to stdout, which is suitable for containerised deployments.
package watcher
