package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/internal/watcher"
)

// historyProvider is satisfied by *watcher.History.
type historyProvider interface {
	Snapshot() []watcher.Record
}

// historyRecord is the JSON-serialisable view of a watcher.Record.
type historyRecord struct {
	JobName   string    `json:"job_name"`
	Kind      string    `json:"kind"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryHandler returns the recent event history as JSON.
//
//	GET /history
func HistoryHandler(h historyProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		records := h.Snapshot()
		out := make([]historyRecord, 0, len(records))
		for _, rec := range records {
			out = append(out, historyRecord{
				JobName:   rec.JobName,
				Kind:      rec.Kind,
				Message:   rec.Message,
				Timestamp: rec.Timestamp,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
