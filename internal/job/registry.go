package job

import (
	"fmt"

	"github.com/cronwatch/internal/config"
)

// Registry holds all monitored jobs indexed by name.
type Registry struct {
	jobs map[string]*Job
}

// NewRegistry builds a Registry from the application config.
func NewRegistry(cfg *config.Config) (*Registry, error) {
	r := &Registry{jobs: make(map[string]*Job, len(cfg.Jobs))}
	for _, jcfg := range cfg.Jobs {
		if _, exists := r.jobs[jcfg.Name]; exists {
			return nil, fmt.Errorf("duplicate job name: %q", jcfg.Name)
		}
		j, err := NewJob(jcfg)
		if err != nil {
			return nil, err
		}
		r.jobs[jcfg.Name] = j
	}
	return r, nil
}

// Get returns the job with the given name, or an error if not found.
func (r *Registry) Get(name string) (*Job, error) {
	j, ok := r.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job not found: %q", name)
	}
	return j, nil
}

// All returns a slice of all registered jobs.
func (r *Registry) All() []*Job {
	out := make([]*Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		out = append(out, j)
	}
	return out
}

// Len returns the number of registered jobs.
func (r *Registry) Len() int { return len(r.jobs) }
