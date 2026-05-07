package watcher

import (
	"os"
	"strings"
	"testing"
)

func TestSyslogWriter_WriteEvent(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "syslog-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	w := NewSyslogWriter(f)
	if err := w.WriteEvent("WARN", "backup-job", "missed run"); err != nil {
		t.Fatalf("WriteEvent: %v", err)
	}

	_ = f.Sync()
	data, _ := os.ReadFile(f.Name())
	line := string(data)

	if !strings.Contains(line, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", line)
	}
	if !strings.Contains(line, "job=backup-job") {
		t.Errorf("expected job=backup-job in output, got: %s", line)
	}
	if !strings.Contains(line, `msg="missed run"`) {
		t.Errorf("expected msg=\"missed run\" in output, got: %s", line)
	}
}

func TestSyslogWriter_WriteStart(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "syslog-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	w := NewSyslogWriter(f)
	if err := w.WriteStart(); err != nil {
		t.Fatalf("WriteStart: %v", err)
	}

	_ = f.Sync()
	data, _ := os.ReadFile(f.Name())
	if !strings.Contains(string(data), "watcher started") {
		t.Errorf("expected 'watcher started', got: %s", string(data))
	}
}

func TestSyslogWriter_WriteStop(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "syslog-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	w := NewSyslogWriter(f)
	if err := w.WriteStop(); err != nil {
		t.Fatalf("WriteStop: %v", err)
	}

	_ = f.Sync()
	data, _ := os.ReadFile(f.Name())
	if !strings.Contains(string(data), "watcher stopped") {
		t.Errorf("expected 'watcher stopped', got: %s", string(data))
	}
}

func TestSyslogWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when nil is passed
	w := NewSyslogWriter(nil)
	if w.out != os.Stdout {
		t.Errorf("expected stdout as default output")
	}
}
