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
