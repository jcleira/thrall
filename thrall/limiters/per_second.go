package limiters

import (
	"sync"
	"time"
)

// PerSecond struct contains all the necessary configuration to setup a Jobs per
// Second limit on the Workers. It have been designed to avoid throttling on the
// Twitter API.
type PerSecond struct {
	Max      int
	Started  int
	Finished int
	Starts   time.Time
	Ends     time.Time
	sync.Mutex
}

// Init initialize the PerSecond limiter. It creates a go routing that would
// reset the Finished jobs every Second.
//
// Returns nothing.
func (ps *PerSecond) Init() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			ps.Lock()
			ps.Finished = 0
			ps.Unlock()
		}
	}()
}

// Adquire checks and adquire a Job if the number of already Started plus
// already Finished Jobs is lesser than the defined Max jobs.
//
// Returns true if the adquire was succesfull, false otherwise.
func (ps *PerSecond) Adquire() bool {
	ps.Lock()
	defer ps.Unlock()

	if ps.Started+ps.Finished >= ps.Max {
		return false
	}

	ps.Started++
	return true
}

// Release releases a started job by converting it to finished, finished jobs
// are cleaned every second.
//
// Returns nothing.
func (ps *PerSecond) Release() {
	ps.Lock()
	defer ps.Unlock()
	ps.Started--
	ps.Finished++
}
