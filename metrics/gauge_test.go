package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewGauges(t *testing.T) {
	assert := assert.New(t)

	t.Run("when NewGauges succeed creating a gauge", func(t *testing.T) {
		r := &Registry{
			Gauges: make(map[string]prometheus.Gauge),
		}

		err := r.NewGauges("foo")

		assert.Nil(err)
		assert.Equal(len(r.Gauges), 1)
		assert.Contains(r.Gauges, "foo")

		err = r.CloseGauge("foo")
		assert.Nil(err)
		assert.Equal(len(r.Gauges), 0)
	})

	t.Run("when NewGauges fails creating a gauge", func(t *testing.T) {
		t.Run("due empty name", func(t *testing.T) {
			r := &Registry{
				Gauges: make(map[string]prometheus.Gauge),
			}

			err := r.NewGauges("")

			assert.Equal("gauge's name should not be empty", err.Error())
			assert.Empty(r.Gauges)
		})

		t.Run("due re-creating the gauge twice", func(t *testing.T) {
			r := &Registry{
				Gauges: make(map[string]prometheus.Gauge),
			}

			err := r.NewGauges("foo")

			assert.Nil(err)

			err = r.NewGauges("foo")

			assert.Equal("gauge 'foo' already registered", err.Error())
			assert.Equal(len(r.Gauges), 1)
			assert.Contains(r.Gauges, "foo")

			err = r.CloseGauge("foo")
			assert.Nil(err)
			assert.Equal(len(r.Gauges), 0)
		})
	})
}

func TestCloseGauge(t *testing.T) {
	assert := assert.New(t)

	t.Run("when CloseGauge succeed on closing a gauge", func(t *testing.T) {
		r := &Registry{
			Gauges: make(map[string]prometheus.Gauge),
		}

		err := r.NewGauges("foo")

		assert.Nil(err)

		assert.Equal(len(r.Gauges), 1)
		assert.Contains(r.Gauges, "foo")

		err = r.CloseGauge("foo")
		assert.Nil(err)
		assert.Equal(len(r.Gauges), 0)
	})

	t.Run("when CloseGauge fails on closing a gauge because it doesn't exists", func(t *testing.T) {
		r := &Registry{
			Gauges: make(map[string]prometheus.Gauge),
		}

		err := r.CloseGauge("foo")
		assert.Equal("gauge 'foo' not registered", err.Error())
		assert.Equal(len(r.Gauges), 0)
	})
}
