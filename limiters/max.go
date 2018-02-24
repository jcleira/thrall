package limiters

import "sync"

// Max struct contains all the necessary configuration to setup a global Max
// jobs limit on the Workers, It's not used right now and has been created by
// error while trying to create the PerSecond limiter.
type Max struct {
	Max  int
	Busy int
	sync.Mutex
}

// Init does nothing but it's necesary to implement the Limiter interface.
//
// Returns nothing.
func (m *Max) Init() {}

// Adquire checks and adquire a Job if the number of currently running jobs or
// Busy jobs is lesser than the defined Max Jobs.
//
// Returns true if the adquire was succesfull, false otherwise.
func (m *Max) Adquire() bool {
	m.Lock()
	defer m.Unlock()

	if m.Busy >= m.Max {
		return false
	}

	m.Busy++
	return true
}

// Release releases a Busy job.
//
// Returns nothing.
func (m *Max) Release() {
	m.Lock()
	defer m.Unlock()
	m.Busy--
}
