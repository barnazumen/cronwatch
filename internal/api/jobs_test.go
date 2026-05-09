package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/job"
)

func makeJobRegistry(t *testing.T) *job.Registry {
	t.Helper()
	cfg := &config.Config{
		Jobs: []config.JobConfig{
			{Name: "backup", Schedule: "@daily", Timeout: "5m"},
			{Name: "cleanup", Schedule: "@hourly", Timeout: "1m"},
		},
	}
	reg, err := job.NewRegistry(cfg)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	return reg
}

func TestJobsHandler_ReturnsAllJobs(t *testing.T) {
	reg := makeJobRegistry(t)
	h := api.JobsHandler(reg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var summaries []api.JobSummary
	if err := json.NewDecoder(rec.Body).Decode(&summaries); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(summaries))
	}
}

func TestJobsHandler_ContentType(t *testing.T) {
	reg := makeJobRegistry(t)
	h := api.JobsHandler(reg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	h(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestJobsHandler_StatusReflectsState(t *testing.T) {
	reg := makeJobRegistry(t)
	j, _ := reg.Get("backup")
	j.RecordSuccess(time.Now())

	h := api.JobsHandler(reg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	h(rec, req)

	var summaries []api.JobSummary
	_ = json.NewDecoder(rec.Body).Decode(&summaries)

	for _, s := range summaries {
		if s.Name == "backup" && s.Status != "ok" {
			t.Errorf("expected status ok, got %q", s.Status)
		}
	}
}

func TestJobsHandler_MethodNotAllowed(t *testing.T) {
	reg := makeJobRegistry(t)
	h := api.JobsHandler(reg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", nil)
	h(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
