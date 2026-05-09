package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/job"
	"github.com/cronwatch/internal/watcher"
)

func makeRouter(t *testing.T) http.Handler {
	t.Helper()
	cfg := &config.Config{
		Jobs: []config.JobConfig{
			{Name: "sync", Schedule: "@hourly", Timeout: "2m"},
		},
	}
	reg, err := job.NewRegistry(cfg)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	metrics := watcher.NewMetrics()
	hc := watcher.NewHealthCheck()
	return api.NewRouter(api.RouterConfig{
		Registry:    reg,
		Metrics:     metrics,
		HealthCheck: hc,
	})
}

func TestRouter_HealthzRoute(t *testing.T) {
	router := makeRouter(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("/healthz: expected 200, got %d", rec.Code)
	}
}

func TestRouter_MetricsRoute(t *testing.T) {
	router := makeRouter(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("/metrics: expected 200, got %d", rec.Code)
	}
}

func TestRouter_JobsRoute(t *testing.T) {
	router := makeRouter(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("/api/jobs: expected 200, got %d", rec.Code)
	}
}

func TestRouter_UnknownRoute(t *testing.T) {
	router := makeRouter(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown route: expected 404, got %d", rec.Code)
	}
}
