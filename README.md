# Thrall

[![circle-ci](https://circleci.com/gh/jcleira/thrall/tree/dev.png?style=shield)](https://circleci.com/gh/jcleira/thrall)

Thrall is a package that provides job processing in Go.

![Thrall photo](https://i.imgur.com/vlMYoCv.jpg)

## Usage

thrall provides a unique public method `Init` that would return two channels, a `chan Runnable` jobs channel and a `chan boolean` quit channel.

`Runnable` is thrall's Go interface that defines a job that could be run, 
```go
// Runnable is thrall's Job interface, that should be implemented by the
// package's users to create jobs tha thrall would be able to run.
type Runnable interface {
	Run() error
}
```

To use thrall create your Jobs structs that implement the Runnable interface `Run() error` function. Then you would be able to send them to the `chan Runnable` returned by the `Init` function.

```go
type Job struct {
  Executed bool 
}

func (j *Job) Run() error {
  j.Executed = true
  return nil
}

jobs, quit := thrall.Init(1) // Initialize thrall with one worker
jobs <- &Job{} 
quit <- true

fmt.Println(job.Executed)
//Output: true
```

## Limiters 

Thrall also provides a set of concurrent jobs limiters that would be useful to control it's deployment or throttling use cases.

Initialize thrall with 8 workers and with 1000 concurrent jobs limit.
```go
jobs, quit := thrall.Init(8, WithMaxLimiter(1000))
```

Initialize thrall with 8 workers and with 16 per second jobs limit.
```go
jobs, quit := thrall.Init(8, WithPersecondLimiter(16))
```
