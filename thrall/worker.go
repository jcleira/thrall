package thrall

import (
	"log"
	"time"
)

// jobTimeout is the defined max time for a job to be executed on the worker.
const jobTimeout = 5 * time.Second

// worker defines a job runner, it belongs to a workerPool and contains a
// reference of it. This abstraction has been created to "control" the number
// thrall's goroutine spawn, we could for each job create a goroutine, but
// instead we run the jobs by attaching them to a limited number of workers.
type worker struct {
	Id         int
	workerPool *workerPool
}

// Worker Start starts the worker goroutine and becomes ready to accept Jobs.
//
// Returns nothing.
func (w *worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.workerPool.Jobs:
				w.Enqueue(job)
			case <-w.workerPool.Quit:
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
		go func() {
			w.workerPool.Jobs <- job
		}()

		return
	}

	w.Run(job)
	w.workerPool.Limiter.Release()
}

// Run executes the Runnable and gives it a constant time to finish, if that
// time is hit, Run free the worker to accept more Jobs but as it launches a
// new goroutine to execute the Runnable that goroutine may get leaked.
//
// TODO a library may not write directly to the log.
//
// - job: The Runnable to run on the worker.
//
// Returns nothing.
func (w *worker) Run(job Runnable) {
	done := make(chan bool)
	go func() {
		if err := job.Run(); err != nil {
			log.Printf("worker %d - job %d - %v", err)
		}

		done <- true
	}()

	select {
	case <-time.After(jobTimeout):
		log.Printf("job timeout (%d sec) on worker %d", w.Id, jobTimeout)
	case <-done:
		return
	}
}
