package limiters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxAquire(t *testing.T) {
	assert := assert.New(t)

	t.Run("when max limiter adquire succeed", func(t *testing.T) {
		max := Max{
			Max:  9,
			Busy: 8,
		}

		result := max.Adquire()
		assert.True(result)
		assert.Equal(9, max.Busy)
	})

	t.Run("when max limiter adquire fails", func(t *testing.T) {
		max := Max{
			Max:  0,
			Busy: 0,
		}

		result := max.Adquire()
		assert.False(result)
		assert.Equal(0, max.Busy)
	})
}
func TestMaxRelease(t *testing.T) {
	assert := assert.New(t)

	t.Run("when max limiter release succeed", func(t *testing.T) {
		max := Max{
			Max:  0,
			Busy: 1,
		}

		max.Release()
		assert.Equal(0, max.Busy)
	})
}
