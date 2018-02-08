package limiters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerSecondInit(t *testing.T) {
	assert := assert.New(t)

	t.Run("when init resets Finished successfully", func(t *testing.T) {
		perSecond := PerSecond{
			Finished: 2,
		}

		perSecond.Init()
		time.Sleep(1 * time.Second)

		assert.Equal(0, perSecond.Finished)
	})
}

func TestPerSecondAquire(t *testing.T) {
	assert := assert.New(t)

	t.Run("when per second limiter adquire succeed", func(t *testing.T) {
		perSecond := PerSecond{
			Max:      5,
			Started:  2,
			Finished: 2,
		}

		result := perSecond.Adquire()
		assert.True(result)
		assert.Equal(3, perSecond.Started)
		assert.Equal(2, perSecond.Finished)
		assert.Equal(5, perSecond.Max)
	})

	t.Run("when per second limiter adquire fails", func(t *testing.T) {
		perSecond := PerSecond{
			Max:      2,
			Started:  1,
			Finished: 1,
		}

		result := perSecond.Adquire()
		assert.False(result)
		assert.Equal(1, perSecond.Started)
		assert.Equal(1, perSecond.Finished)
		assert.Equal(2, perSecond.Max)
	})

}

func TestPerSecondRelease(t *testing.T) {
	assert := assert.New(t)

	t.Run("when per second limiter release succeed", func(t *testing.T) {
		perSecond := PerSecond{
			Max:      0,
			Started:  1,
			Finished: 1,
		}

		perSecond.Release()
		assert.Equal(0, perSecond.Started)
		assert.Equal(2, perSecond.Finished)
		assert.Equal(0, perSecond.Max)
	})
}
