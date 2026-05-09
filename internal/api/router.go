package api

import (
	"net/http"

	"github.com/cronwatch/internal/job"
	"github.com/cronwatch/internal/watcher"
)

// RouterConfig holds the dependencies needed to build the HTTP router.
type RouterConfig struct {
	Registry    *job.Registry
	Metrics     *watcher.Metrics
	HealthCheck *watcher.HealthCheck
}

// NewRouter constructs and returns the HTTP mux for the cronwatch API.
// Routes:
//
//	GET /healthz      — liveness / readiness probe
//	GET /metrics      — internal counters snapshot
//	GET /api/jobs     — list all monitored jobs and their status
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", HealthzHandler(cfg.Registry, cfg.HealthCheck))
	mux.HandleFunc("/metrics", MetricsHandler(cfg.Metrics))
	mux.HandleFunc("/api/jobs", JobsHandler(cfg.Registry))

	return mux
}
