package thrall

// Runnable is thrall's Job interface, that should be implemented by the
// package's users to create jobs tha thrall would be able to run.
type Runnable interface {
	Run() error
}
