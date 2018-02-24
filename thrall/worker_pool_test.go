package thrall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJob struct {
	Executed bool
}

func (tj *testJob) Run() error {
	tj.Executed = true
	return nil
}

type scheduleableJob struct {
	testJob
}

func (sj *scheduleableJob) Schedule() time.Time {
	return time.Now().Add(1 * time.Hour)
}

type repeatableJob struct {
	testJob
	Repeated int
}

func (r *repeatableJob) Repeat() bool {
	if r.Repeated == 2 {
		return false
	}
	r.Repeated += 1

	return true
}

func TestInit(t *testing.T) {
	assert := assert.New(t)

	t.Run("when Init succeed and initializing thrall", func(t *testing.T) {
		queue, close := Init(1)

		assert.NotNil(queue)
		assert.NotNil(close)
		assert.Len(wp.workers, 1)
		assert.Equal(wp.Queue, queue)
		assert.Equal(wp.close, close)

		close <- true
	})
}

func TestSchedule(t *testing.T) {
	assert := assert.New(t)

	t.Run("when schedule succeed at scheduling a job", func(t *testing.T) {
		queue, close := Init(1)

		queue <- &scheduleableJob{}
		time.Sleep(10 * time.Millisecond)

		assert.Len(wp.Delayed, 1)

		close <- true
	})
}
