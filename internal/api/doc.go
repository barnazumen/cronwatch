// Package api exposes a lightweight HTTP interface for cronwatch.
//
// # Routes
//
//	GET /healthz
//	  Returns 200 OK with a JSON body describing the health of each monitored
//	  job. Returns 503 if any job is in a degraded state.
//
//	GET /metrics
//	  Returns a JSON snapshot of internal watcher counters: total checks
//	  performed, missed jobs detected, and alerts dispatched.
//
//	GET /api/jobs
//	  Returns a JSON array of every registered job together with its current
//	  status, schedule, and last-run timestamps.
//
// # Usage
//
//	router := api.NewRouter(api.RouterConfig{
//	    Registry:    reg,
//	    Metrics:     metrics,
//	    HealthCheck: hc,
//	})
//	http.ListenAndServe(":8080", router)
package api
