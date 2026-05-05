package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
log_level: debug
alerts:
  email: ops@example.com
jobs:
  - name: backup
    schedule: "0 2 * * *"
    command: /usr/local/bin/backup.sh
    timeout: 30m
`
	path := writeTempConfig(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log_level=debug, got %q", cfg.LogLevel)
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "backup" {
		t.Errorf("expected job name=backup, got %q", cfg.Jobs[0].Name)
	}
	if cfg.Jobs[0].Timeout != 30*time.Minute {
		t.Errorf("expected timeout=30m, got %v", cfg.Jobs[0].Timeout)
	}
}

func TestLoad_DefaultLogLevel(t *testing.T) {
	yaml := `
jobs:
  - name: ping
    schedule: "* * * * *"
    command: /bin/ping
`
	path := writeTempConfig(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level=info, got %q", cfg.LogLevel)
	}
}

func TestLoad_MissingJobName(t *testing.T) {
	yaml := `
jobs:
  - schedule: "* * * * *"
    command: /bin/true
`
	path := writeTempConfig(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing job name")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	yaml := `log_level: info
`
	path := writeTempConfig(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when no jobs defined")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
