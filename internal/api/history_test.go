package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/watcher"
)

// stubHistory implements historyProvider for tests.
type stubHistory struct {
	records []watcher.Record
}

func (s *stubHistory) Snapshot() []watcher.Record { return s.records }

func makeHistory(records []watcher.Record) *stubHistory {
	return &stubHistory{records: records}
}

func TestHistoryHandler_ReturnsRecords(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	h := makeHistory([]watcher.Record{
		{JobName: "backup", Kind: "missed", Message: "backup missed", Timestamp: now},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	api.HistoryHandler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var out []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 record, got %d", len(out))
	}
	if out[0]["job_name"] != "backup" {
		t.Errorf("expected job_name 'backup', got %v", out[0]["job_name"])
	}
	if out[0]["kind"] != "missed" {
		t.Errorf("expected kind 'missed', got %v", out[0]["kind"])
	}
}

func TestHistoryHandler_EmptyList(t *testing.T) {
	h := makeHistory(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	api.HistoryHandler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var out []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d items", len(out))
	}
}

func TestHistoryHandler_ContentType(t *testing.T) {
	h := makeHistory(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	api.HistoryHandler(h).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHistoryHandler_MethodNotAllowed(t *testing.T) {
	h := makeHistory(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	api.HistoryHandler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
