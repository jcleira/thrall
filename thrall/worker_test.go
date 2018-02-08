package thrall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJob struct {
	Executed bool
}

func (tj *testJob) Run() error {
	tj.Executed = true
	return nil
}

func TestWorkerEnqueue(t *testing.T) {
	assert := assert.New(t)

	t.Run("when there is no limiters", func(t *testing.T) {
		t.Run("enqueue succeed adding the Job", func(t *testing.T) {
			jobs, quit := Init(1)

			var job testJob
			jobs <- &job
			quit <- true

			assert.True(job.Executed)
		})
	})

	t.Run("when there is limiters", func(t *testing.T) {
		t.Run("enqueue succeed if limit is not reached", func(t *testing.T) {
			jobs, quit := Init(1, WithMaxLimiter(1))

			var job testJob
			jobs <- &job
			quit <- true

			assert.True(job.Executed)
		})

		t.Run("enqueue fails if limit is reached", func(t *testing.T) {
			jobs, quit := Init(1, WithMaxLimiter(0))

			var job testJob
			jobs <- &job
			quit <- true

			assert.False(job.Executed)
		})
	})
}
