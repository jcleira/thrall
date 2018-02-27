package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	assert := assert.New(t)

	t.Run("when NewCounter succeed creating a counter", func(t *testing.T) {
		r := &Registry{
			Counters: make(map[string]prometheus.Counter),
		}

		err := r.NewCounter("foo", "bar")

		assert.Nil(err)
		assert.Equal(len(r.Counters), 1)
		assert.Contains(r.Counters, "foo")

		err = r.CloseCounter("foo")
		assert.Nil(err)
		assert.Equal(len(r.Counters), 0)
	})

	t.Run("when NewCounter fails creating a counter", func(t *testing.T) {
		t.Run("due empty name", func(t *testing.T) {
			r := &Registry{
				Counters: make(map[string]prometheus.Counter),
			}

			err := r.NewCounter("", "bar")

			assert.Equal("counter's name should not be empty", err.Error())
			assert.Empty(r.Counters)
		})

		t.Run("due empty help", func(t *testing.T) {
			r := &Registry{
				Counters: make(map[string]prometheus.Counter),
			}

			err := r.NewCounter("foo", "")

			assert.Equal("counter's 'foo' help should not be empty", err.Error())
			assert.Empty(r.Counters)
		})

		t.Run("due re-creating the counter twice", func(t *testing.T) {
			r := &Registry{
				Counters: make(map[string]prometheus.Counter),
			}

			err := r.NewCounter("foo", "bar")

			assert.Nil(err)

			err = r.NewCounter("foo", "bar")

			assert.Equal("counter 'foo' already registered", err.Error())
			assert.Equal(len(r.Counters), 1)
			assert.Contains(r.Counters, "foo")

			err = r.CloseCounter("foo")
			assert.Nil(err)
			assert.Equal(len(r.Counters), 0)
		})
	})
}
