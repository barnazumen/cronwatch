package job

import (
	"fmt"
	"sync"

	"github.com/cronwatch/cronwatch/internal/config"
)

// Registry holds all monitored jobs indexed by name.
type Registry struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

// NewRegistry builds a Registry from the loaded configuration.
func NewRegistry(cfg *config.Config) (*Registry, error) {
	r := &Registry{
		jobs: make(map[string]*Job, len(cfg.Jobs)),
	}
	for _, jcfg := range cfg.Jobs {
		if _, exists := r.jobs[jcfg.Name]; exists {
			return nil, fmt.Errorf("duplicate job name: %q", jcfg.Name)
		}
		r.jobs[jcfg.Name] = NewJob(jcfg.Name, jcfg.Schedule, jcfg.GracePeriod)
	}
	return r, nil
}

// Get returns the job with the given name, or an error if not found.
func (r *Registry) Get(name string) (*Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job not found: %q", name)
	}
	return j, nil
}

// All returns a slice of snapshots for every registered job.
func (r *Registry) All() []Job {
	r.mu.RLock()
	defer r.mu.RUnlock()
	snaps := make([]Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		snaps = append(snaps, j.Snapshot())
	}
	return snaps
}

// Len returns the number of registered jobs.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.jobs)
}
