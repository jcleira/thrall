package limiters

type Limiter interface {
	Init()
	Adquire() bool
	Release()
}
