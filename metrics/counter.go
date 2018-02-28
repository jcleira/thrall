package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCounter creates N number of new Counter metrics.
//
// - names: Counter's names.
//
// Returns an error if any counter creation fails.
func (r *Registry) NewCounters(names ...string) error {
	r.Lock()
	defer r.Unlock()

	for _, name := range names {
		if name == "" {
			return errors.New("counter's name should not be empty")
		}

		if _, exists := r.Counters[name]; exists {
			return fmt.Errorf("counter '%s' already registered", name)
		}

		r.Counters[name] = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name,
				Help: name,
			},
		)

		prometheus.MustRegister(r.Counters[name])
	}

	return nil
}

// CloseCounter unregister and remove an already created Counter.
//
// - name: Counter's name to close.
//
// Returns an error if the counter is not found.
func (r *Registry) CloseCounter(name string) error {
	if _, exists := r.Counters[name]; !exists {
		return fmt.Errorf("counter '%s' not registered", name)
	}

	if unregistered := prometheus.Unregister(r.Counters[name]); !unregistered {
		return fmt.Errorf("counter '%s' not unregistered", name)
	}

	r.Lock()
	defer r.Unlock()

	delete(r.Counters, name)

	return nil
}
