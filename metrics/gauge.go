package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// NewGauge creates a new Gauge metric.
//
// - name: Gauge's name.
// - help: Gauge's help, required by prometheus.
//
// Returns nothing.
func (r *Registry) NewGauge(name, help string) error {
	if name == "" {
		return errors.New("gauge's name should not be empty")
	}

	if help == "" {
		return fmt.Errorf("gauge's '%s' help should not be empty", name)
	}

	if _, exists := r.Gauges[name]; exists {
		return fmt.Errorf("gauge '%s' already registered", name)
	}

	r.Lock()
	defer r.Unlock()

	r.Gauges[name] = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
	)

	prometheus.Register(r.Gauges[name])

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

// IncGauge increments by one a gauge's value.
//
// - name: gauge's name to increment.
//
// Returns an error if the gauge is not found.
func (r *Registry) IncGauge(name string) error {
	gauge, exists := r.Gauges[name]
	if !exists {
		return fmt.Errorf("gauge %s not found", name)
	}

	gauge.Inc()

	return nil
}

// IncGauge decrements by one a gauge's value.
//
// - name: gauge's name to decrement.
//
// Returns an error if the gauge is not found.
func (r *Registry) DecGauge(name string) error {
	gauge, exists := r.Gauges[name]
	if !exists {
		return fmt.Errorf("gauge %s not found", name)
	}

	gauge.Inc()

	return nil
}

// SetGauge sets Gauge's value.
//
// - name: Gauge's name to set it's value.
//
// Returns an error if the gauge is not found.
func (r *Registry) SetGauge(name string, value float64) error {
	gauge, exists := r.Gauges[name]
	if !exists {
		return fmt.Errorf("gauge %s not found", name)
	}

	gauge.Set(value)

	return nil
}
