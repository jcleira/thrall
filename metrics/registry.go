package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry is metrics main struct it's a metrics container. To use the metrics
// package you should first create a Registry, then start adding metrics to the
// Registry, then modify those metrics values.
type Registry struct {
	Gauges   map[string]prometheus.Gauge
	Counters map[string]prometheus.Counter

	sync.Mutex
}

// NewRegistry creates an empty registry and configures the metrics endpoint.
//
// Returns an empty Registry.
func NewRegistry() *Registry {
	http.Handle("/metrics", promhttp.Handler())

	return &Registry{
		Gauges:   make(map[string]prometheus.Gauge),
		Counters: make(map[string]prometheus.Counter),
	}
}

// Inc increases the value of N number of metrics given by it's names.
//
// - names: The metric names to increase.
//
// Returns nothing.
func (r *Registry) Inc(names ...string) {
	for _, name := range names {
		if gauge, exists := r.Gauges[name]; exists {
			gauge.Inc()
		}

		if counter, exists := r.Counters[name]; exists {
			counter.Inc()
		}
	}
}

// Dec decreases the value of N number of metrics given by it's names.
//
// - names: The metric names to decrease.
//
// Returns nothing.
func (r *Registry) Dec(names ...string) {
	for _, name := range names {
		if gauge, exists := r.Gauges[name]; exists {
			gauge.Dec()
		}
	}
}
