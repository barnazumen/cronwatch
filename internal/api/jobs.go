package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/internal/job"
)

// JobSummary is the JSON representation of a single job's current state.
type JobSummary struct {
	Name        string    `json:"name"`
	Schedule    string    `json:"schedule"`
	Status      string    `json:"status"`
	LastRun     time.Time `json:"last_run,omitempty"`
	LastSuccess time.Time `json:"last_success,omitempty"`
	LastFailure time.Time `json:"last_failure,omitempty"`
	MissedCount int       `json:"missed_count"`
}

// JobsHandler returns a JSON array of all registered jobs and their current status.
func JobsHandler(reg *job.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		names := reg.Names()
		summaries := make([]JobSummary, 0, len(names))

		for _, name := range names {
			j, ok := reg.Get(name)
			if !ok {
				continue
			}
			summaries = append(summaries, JobSummary{
				Name:        j.Name,
				Schedule:    j.Schedule,
				Status:      j.Status.String(),
				LastRun:     j.LastRun,
				LastSuccess: j.LastSuccess,
				LastFailure: j.LastFailure,
				MissedCount: j.MissedCount,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(summaries)
	}
}
