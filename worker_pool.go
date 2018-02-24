package thrall

import (
	"sync"
	"time"

	"github.com/jcleira/thrall/limiters"
)

// workerPool defines a group of workers and their common characteristics. The
// Jobs Runnable channel would feed all the workers on the pool with jobs. The
// close channel would make all the workers to finish whne called, and the
// Limiters would set a group for rules for all the workers on the pool to
// follow.
type workerPool struct {
	// Queue is the main Jobs Queue, It contains all the jobs that are waiting to
	// be run.
	Queue chan Runnable

	// Delayed is the Jobs Scheduled Queue, It containes all the scheduled jobs
	// that are waiting for their execution time.
	Delayed       map[time.Time][]Runnable
	DelayedMutext sync.Mutex

	// Limiter si thrall configured Jobs limiter, currently only one limiter can
	// be configured per workerPool
	Limiter limiters.Limiter

	workersQueue chan Runnable
	workersClose chan bool
	workers      []*worker
	errors       chan error
	close        chan bool
}

// wp is a package var that will keep all thrall's state.
var wp *workerPool

// Init is thrall's main initializer, it actually initizalices and run a
// workerPool and it's defined workers. Init is also the thrall's public API,
// it does return the two thrall's interactors, a chan Runnable that is the Job
// queue and a quit channel to stop thrall's world.
//
// - numWorkers: the number of workers for the thrall's workPool.
// - opts: function option initializers, check the following With.. funcs.
//
// Returns the jobs queue, an errors channel and close channel.
func Init(workers int, opts ...func(*workerPool)) (chan Runnable, chan error, chan bool) {
	wp = &workerPool{
		Queue:        make(chan Runnable),
		Delayed:      make(map[time.Time][]Runnable),
		close:        make(chan bool),
		errors:       make(chan error),
		workersQueue: make(chan Runnable),
		workersClose: make(chan bool),
	}

	for _, option := range opts {
		option(wp)
	}

	// TODO we are forcing one limiter, we might force the user to send it.
	if len(opts) == 0 {
		wp.Limiter = &limiters.Max{Max: 1000}
	}

	for i := 1; i <= workers; i++ {
		worker := &worker{
			Id:         i,
			workerPool: wp,
			Queue:      wp.workersQueue,
			Errors:     wp.errors,
			Close:      wp.workersClose,
		}

		wp.workers = append(wp.workers, worker)
	}

	wp.run()

	return wp.Queue, wp.errors, wp.close
}

// WithMaxLimiter is an optional func for thrall's init, It does configure a
// max concurrent job limiter for thrall, said otherwise, all the workers
// won't execute more jobs than the maxJobs number given concurrently.
//
// - maxJobs: The max number of concurrent jobs for thrall.
//
// Returns a optional configuration function.
func WithMaxLimiter(maxJobs int) func(*workerPool) {
	return func(wp *workerPool) {
		wp.Limiter = &limiters.Max{Max: maxJobs}
	}
}

// WithPerSecondLimiter is an optional func for thrall's init, It does configure
// a max concurrent job per second limiter for thrall, within all the workers no
// more than the given perSecondJobs param would be executed on the same second.
//
// - perSecondJobs: The max number of jobs jobs per second for thrall.
//
// Returns a optional configuration function.
func WithPerSecondLimiter(perSecondJobs int) func(*workerPool) {
	return func(wp *workerPool) {
		wp.Limiter = &limiters.PerSecond{Max: perSecondJobs}
	}
}

// run would launch the workerPool by starting all it's workers.
//
// Returns nothing.
func (wp *workerPool) run() {
	go func() {
		for {
			select {
			case job := <-wp.Queue:
				if scheduleable, ok := job.(Scheduleable); ok {
					go wp.schedule(job, scheduleable.Schedule())
					continue
				}

				wp.workersQueue <- job
			case <-wp.close:
				close(wp.workersClose)
				return
			}
		}
	}()

	for i := 0; i < len(wp.workers); i++ {
		wp.workers[i].Start()
	}
}

// schedule performs job scheduling for thrall's scheduleable job interfaces
//
// - job: The job to schedule.
// - when: The job programmed execution time.
//
// Returns nothing.
func (wp *workerPool) schedule(job Runnable, when time.Time) {
	wp.DelayedMutext.Lock()
	defer wp.DelayedMutext.Unlock()

	if delayed, ok := wp.Delayed[when]; ok {
		delayed = append(delayed, job)
	} else {
		wp.Delayed[when] = []Runnable{job}
	}
}
