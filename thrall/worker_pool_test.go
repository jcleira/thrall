package thrall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)

	t.Run("when run with some workers succeed", func(t *testing.T) {
		jobs, quit := Init(1)

		assert.NotNil(jobs)
		assert.NotNil(quit)

		quit <- true
	})
}
