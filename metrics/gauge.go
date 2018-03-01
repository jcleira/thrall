package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// NewGauges creates N number of new Gauge metrics.
//
// - name: Gauge's name.
//
// Returns an error if any gauge creation fails.
func (r *Registry) NewGauges(names ...string) error {
	r.Lock()
	defer r.Unlock()

	for _, name := range names {
		if name == "" {
			return errors.New("gauge's name should not be empty")
		}

		if _, exists := r.Gauges[name]; exists {
			return fmt.Errorf("gauge '%s' already registered", name)
		}

		r.Gauges[name] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: name,
				Help: name,
			},
		)

		prometheus.Register(r.Gauges[name])
	}

	return nil
}

// CloseGauge unregister and remove an already created Gauge.
//
// - name: Gauge's name to close.
//
// Returns an error if the gauge is not found.
func (r *Registry) CloseGauge(name string) error {
	if _, exists := r.Gauges[name]; !exists {
		return fmt.Errorf("gauge '%s' not registered", name)
	}

	if unregistered := prometheus.Unregister(r.Gauges[name]); !unregistered {
		return fmt.Errorf("gauge '%s' not unregistered", name)
	}

	r.Lock()
	defer r.Unlock()

	delete(r.Gauges, name)

	return nil
}
