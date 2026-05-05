package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/cronwatch/internal/job"
)

// LogAlerter is a simple Alerter implementation that writes alerts to the
// standard logger. It is useful as a default when no external notification
// channel is configured.
type LogAlerter struct {
	prefix string
}

// NewLogAlerter returns a LogAlerter that prefixes every message with prefix.
func NewLogAlerter(prefix string) *LogAlerter {
	if prefix == "" {
		prefix = "ALERT"
	}
	return &LogAlerter{prefix: prefix}
}

// Alert logs the alert details and returns nil.
func (l *LogAlerter) Alert(jobName string, status job.Status, message string) error {
	log.Printf(
		"[%s] %s | job=%s status=%s message=%s",
		l.prefix,
		time.Now().UTC().Format(time.RFC3339),
		jobName,
		status,
		message,
	)
	return nil
}

// FormatAlert returns a human-readable alert string without logging it.
func FormatAlert(jobName string, status job.Status, message string) string {
	return fmt.Sprintf(
		"job=%s status=%s message=%s at=%s",
		jobName,
		status,
		message,
		time.Now().UTC().Format(time.RFC3339),
	)
}
