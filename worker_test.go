package thrall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerEnqueue(t *testing.T) {
	assert := assert.New(t)

	t.Run("when there is no limiters", func(t *testing.T) {
		t.Run("enqueue succeed adding the Job", func(t *testing.T) {
			queue, _, close := Init(1)

			var job testJob
			queue <- &job
			time.Sleep(10 * time.Millisecond)

			assert.True(job.Executed)

			close <- true
		})
	})

	t.Run("when there is limiters", func(t *testing.T) {
		t.Run("enqueue succeed if limit is not reached", func(t *testing.T) {
			queue, _, close := Init(1, WithMaxLimiter(1))

			var job testJob
			queue <- &job
			time.Sleep(10 * time.Millisecond)

			assert.True(job.Executed)

			close <- true
		})

		t.Run("enqueue fails if limit is reached", func(t *testing.T) {
			queue, _, close := Init(1, WithMaxLimiter(0))

			var job testJob
			queue <- &job
			time.Sleep(10 * time.Millisecond)

			assert.False(job.Executed)

			close <- true
		})
	})
}
func TestWorkerRun(t *testing.T) {
	assert := assert.New(t)

	t.Run("when a Repeatable Job succeed on getting repeated", func(t *testing.T) {
		queue, _, close := Init(1)

		var job repeatableJob
		queue <- &job
		time.Sleep(10 * time.Millisecond)

		assert.True(job.Executed)
		assert.Equal(2, job.Repeated)

		close <- true
	})

	t.Run("When a job fails the error is returned properly", func(t *testing.T) {
		queue, errors, close := Init(1)

		queue <- &errorJob{}
		select {
		case err := <-errors:
			assert.Equal("job error on worker 1. Err: error!", err.Error())
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for the job to return an error")
		}
		close <- true
	})
}
