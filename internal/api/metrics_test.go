package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/cronwatch/internal/api"
	"github.com/cronwatch/cronwatch/internal/watcher"
)

func makeMetrics(checks, missed, alerted int) *watcher.Metrics {
	m := watcher.NewMetrics()
	for i := 0; i < checks; i++ {
		m.RecordCheck()
	}
	for i := 0; i < missed; i++ {
		m.RecordMissed()
	}
	for i := 0; i < alerted; i++ {
		m.RecordAlerted()
	}
	return m
}

func TestMetricsHandler_ReturnsJSON(t *testing.T) {
	m := makeMetrics(10, 2, 1)
	h := api.MetricsHandler(m)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var snap api.MetricsSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if snap.TotalChecks != 10 {
		t.Errorf("TotalChecks: want 10, got %d", snap.TotalChecks)
	}
	if snap.TotalMissed != 2 {
		t.Errorf("TotalMissed: want 2, got %d", snap.TotalMissed)
	}
	if snap.TotalAlerted != 1 {
		t.Errorf("TotalAlerted: want 1, got %d", snap.TotalAlerted)
	}
}

func TestMetricsHandler_ContentType(t *testing.T) {
	m := watcher.NewMetrics()
	h := api.MetricsHandler(m)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type: want application/json, got %s", ct)
	}
}

func TestMetricsHandler_ZeroValues(t *testing.T) {
	m := watcher.NewMetrics()
	h := api.MetricsHandler(m)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h(rec, req)

	var snap api.MetricsSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if snap.TotalChecks != 0 || snap.TotalMissed != 0 || snap.TotalAlerted != 0 {
		t.Errorf("expected all zeros, got %+v", snap)
	}
}
