package thrall

// Repeatable defines an interface that should be implemented for that jobs
// that would need to be re-executed repeatedly, the Repeat() func may be used
// to control the repetition number.
type Repeateable interface {
	Repeat() bool
}
