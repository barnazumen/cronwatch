package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatch/cronwatch/internal/watcher"
)

// MetricsSnapshot is the JSON shape returned by GET /metrics.
type MetricsSnapshot struct {
	TotalChecks  int64 `json:"total_checks"`
	TotalMissed  int64 `json:"total_missed"`
	TotalAlerted int64 `json:"total_alerted"`
}

// MetricsHandler returns an HTTP handler that exposes watcher metrics as JSON.
//
// GET /metrics
//
//	200 OK  { "total_checks": 42, "total_missed": 1, "total_alerted": 1 }
func MetricsHandler(m *watcher.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := m.Snapshot()

		payload := MetricsSnapshot{
			TotalChecks:  snap.TotalChecks,
			TotalMissed:  snap.TotalMissed,
			TotalAlerted: snap.TotalAlerted,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
