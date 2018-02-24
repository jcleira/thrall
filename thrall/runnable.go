package thrall

// Runnable is thrall's main Job interface, the Run() func is some work that
// need to be performed. It should be implemented by the package's users to
// create jobs to be run using thrall.
type Runnable interface {
	Run() error
}
