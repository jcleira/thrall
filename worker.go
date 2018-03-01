package thrall

import (
	"fmt"
	"time"
)

// jobTimeout is the defined max time for a job to be executed on the worker.
const jobTimeout = 30 * time.Second

// worker defines a job runner, it belongs to a workerPool and contains a
// reference of it. This abstraction has been created to "control" the number
// thrall's goroutine spawn, we could for each job create a goroutine, but
// instead we run the jobs by attaching them to a limited number of workers.
type worker struct {
	Id         int
	Queue      chan Runnable
	Errors     chan error
	Close      chan bool
	workerPool *workerPool
}

// Worker Start starts the worker goroutine and becomes ready to accept Jobs.
//
// Returns nothing.
func (w *worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.Queue:
				w.Enqueue(job)
			case <-w.Close:
				return
			}
		}
	}()
}

// Enqueue feeds the worker with a job, but it might don't ingest it under under
// some circunstances as hiting the configured limiter's limit.
//
// - job: The Runnable to enqueue on the worker.
//
// Returns nothing.
func (w *worker) Enqueue(job Runnable) {
	if !w.workerPool.Limiter.Adquire() {
		w.workerPool.IncMetric("thrall_workerpool_job_rate_limited")
		go func() {
			w.workerPool.Queue <- job
		}()

		return
	}

	w.Run(job)
	w.workerPool.Limiter.Release()

	if repeatable, ok := job.(Repeateable); ok {
		if repeatable.Repeat() {
			w.workerPool.Queue <- job
			return
		}
	}
}

// Run executes the Runnable and gives it a constant time to finish, if that
// time is hit, Run free the worker to accept more Jobs but as it launches a
// new goroutine to execute the Runnable that goroutine may get leaked.
//
// - job: The Runnable to run on the worker.
//
// Returns nothing.
func (w *worker) Run(job Runnable) {
	done := make(chan bool)
	go func() {
		if err := job.Run(); err != nil {
			w.workerPool.IncMetric("thrall_workerpool_job_erroed")
			w.Errors <- fmt.Errorf("job error on worker %d. Err: %v", w.Id, err)
		}

		w.workerPool.IncMetric("thrall_workerpool_job_processed")
		done <- true
	}()

	select {
	case <-time.After(jobTimeout):
		w.workerPool.IncMetric("thrall_workerpool_job_timeout")
		w.Errors <- fmt.Errorf("job timeout (%f sec) on worker %d", jobTimeout.Seconds(), w.Id)
	case <-done:
	}

	w.workerPool.DecMetric("thrall_workerpool_job_enqueued")
}
