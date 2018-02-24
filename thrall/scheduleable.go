package thrall

import "time"

// Scheduleable defines an interface that should be implemented for that jobs
// that would need to be executed on given time, the Schedule() func returns
// a time.Time that would be the execution time for that job.
type Scheduleable interface {
	Schedule() time.Time
}
