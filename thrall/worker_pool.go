package thrall

import (
	"github.com/jcleira/thrall/thrall/limiters"
)

// workerPool defines a group of Workers and their common characteristics. The
// Jobs Runnable channel would feed all the workers on the pool with jobs. The
// Quit channel would make all the workers to finish whne called, and the
// Limiters would set a group for rules for all the workers on the pool to
// follow.
type workerPool struct {
	Workers []*worker
	Jobs    chan Runnable
	Quit    chan bool
	Limiter limiters.Limiter
}

// Init is thrall's main initializer, it actually initizalices and run a
// workerPool and it's defined workers. Init is also the thrall's public API,
// it does return the two thrall's interactors, a chan Runnable that is the Job
// queue and a quit channel to stop thrall's world.
//
// - numWorkers: the number of Workers for the thrall's workPool.
// - options: function option initializers, check the following With.. funcs.
//
// Returns the jobs channel as chan Runnable, and a boolean quit channel.
func Init(numWorkers int, options ...func(*workerPool)) (chan Runnable, chan bool) {
	jobs := make(chan Runnable)
	quit := make(chan bool)

	wp := &workerPool{
		Jobs: jobs,
		Quit: quit,
	}

	for _, option := range options {
		option(wp)
	}

	// TODO we are forcing one limiter, we might force the user to send it.
	if len(options) == 0 {
		wp.Limiter = &limiters.Max{Max: 1000}
	}

	for i := 0; i < numWorkers; i++ {
		worker := &worker{
			Id:         i,
			workerPool: wp,
		}

		wp.Workers = append(wp.Workers, worker)
	}

	wp.run()

	return jobs, quit
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

// run would launch the workerPool by starting all it's Workers.
//
// Returns nothing.
func (wp *workerPool) run() {
	for i := 0; i < len(wp.Workers); i++ {
		wp.Workers[i].Start()
	}
}
