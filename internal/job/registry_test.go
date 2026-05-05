package job

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
)

func makeConfig(names ...string) *config.Config {
	cfg := &config.Config{}
	for _, n := range names {
		cfg.Jobs = append(cfg.Jobs, config.JobConfig{
			Name:        n,
			Schedule:    "* * * * *",
			GracePeriod: time.Minute,
		})
	}
	return cfg
}

func TestNewRegistry_OK(t *testing.T) {
	r, err := NewRegistry(makeConfig("jobA", "jobB"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 2 {
		t.Fatalf("expected 2 jobs, got %d", r.Len())
	}
}

func TestNewRegistry_DuplicateName(t *testing.T) {
	_, err := NewRegistry(makeConfig("dup", "dup"))
	if err == nil {
		t.Fatal("expected error for duplicate job name")
	}
}

func TestRegistry_Get(t *testing.T) {
	r, _ := NewRegistry(makeConfig("myJob"))
	j, err := r.Get("myJob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.Name != "myJob" {
		t.Fatalf("expected myJob, got %s", j.Name)
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r, _ := NewRegistry(makeConfig("present"))
	_, err := r.Get("absent")
	if err == nil {
		t.Fatal("expected error for missing job")
	}
}

func TestRegistry_All(t *testing.T) {
	r, _ := NewRegistry(makeConfig("a", "b", "c"))
	all := r.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(all))
	}
}
