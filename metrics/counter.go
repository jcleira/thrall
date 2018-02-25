package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCounter creates a new Counter metric.
//
// - name: Counter's name.
// - help: Counter's help, required by prometheus.
//
// Returns nothing.
func (r *Registry) NewCounter(name, help string) error {
	if name == "" {
		return errors.New("counter's name should not be empty")
	}

	if help == "" {
		return fmt.Errorf("counter's '%s' help should not be empty", name)
	}

	if _, exists := r.Counters[name]; exists {
		return fmt.Errorf("counter '%s' already registered", name)
	}

	r.Lock()
	defer r.Unlock()

	r.Counters[name] = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
	)

	prometheus.MustRegister(r.Counters[name])

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

// IncCounter increments by one a counter's value.
//
// - name: counter's name to increment.
//
// Returns an error if the counter is not found.
func (r *Registry) IncCounter(name string) error {
	counter, exists := r.Counters[name]
	if !exists {
		return fmt.Errorf("counter %s not found", name)
	}

	counter.Inc()

	return nil
}
