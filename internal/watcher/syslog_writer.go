package watcher

import (
	"fmt"
	"os"
	"time"
)

// SyslogWriter writes structured watcher events to a file or stdout
// in a syslog-compatible format for integration with log aggregators.
type SyslogWriter struct {
	out *os.File
}

// NewSyslogWriter creates a SyslogWriter. Pass nil to write to stdout.
func NewSyslogWriter(f *os.File) *SyslogWriter {
	if f == nil {
		f = os.Stdout
	}
	return &SyslogWriter{out: f}
}

// WriteEvent writes a single watcher event line.
func (s *SyslogWriter) WriteEvent(level, jobName, message string) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(s.out, "%s [%s] job=%s msg=%q\n", timestamp, level, jobName, message)
	return err
}

// WriteStart writes a startup event.
func (s *SyslogWriter) WriteStart() error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(s.out, "%s [INFO] cronwatch watcher started\n", timestamp)
	return err
}

// WriteStop writes a shutdown event.
func (s *SyslogWriter) WriteStop() error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(s.out, "%s [INFO] cronwatch watcher stopped\n", timestamp)
	return err
}
