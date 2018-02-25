package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	assert := assert.New(t)

	t.Run("when NewRegistry succeed creating a registry", func(t *testing.T) {
		r := NewRegistry()

		assert.NotNil(r)
		assert.NotNil(r.Gauges)
		assert.NotNil(r.Counters)
	})
}
