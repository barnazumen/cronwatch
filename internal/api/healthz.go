// Package api exposes HTTP endpoints for cronwatch operational status.
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/cronwatch/internal/watcher"
)

// healthzResponse is the JSON shape returned by GET /healthz.
type healthzResponse struct {
	Status string        `json:"status"`
	Jobs   []jobHealth   `json:"jobs"`
	At     time.Time     `json:"at"`
}

type jobHealth struct {
	Name      string    `json:"name"`
	Healthy   bool      `json:"healthy"`
	Message   string    `json:"message"`
	LastCheck time.Time `json:"last_check"`
}

// HealthzHandler returns an http.HandlerFunc that serves health status.
func HealthzHandler(hc *watcher.HealthCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := hc.Snapshot()
		allHealthy := true
		jobs := make([]jobHealth, 0, len(snap))
		for _, s := range snap {
			if !s.Healthy {
				allHealthy = false
			}
			jobs = append(jobs, jobHealth{
				Name:      s.JobName,
				Healthy:   s.Healthy,
				Message:   s.Message,
				LastCheck: s.LastCheck,
			})
		}
		status := "ok"
		httpCode := http.StatusOK
		if !allHealthy {
			status = "degraded"
			httpCode = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		_ = json.NewEncoder(w).Encode(healthzResponse{
			Status: status,
			Jobs:   jobs,
			At:     time.Now().UTC(),
		})
	}
}
