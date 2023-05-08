package utils

type Tuple2[A any, B any] struct {
	Value A
	Err   B
}

func T2[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{Value: a, Err: b}
}

func Async[A any](f func() A) <-chan A {
	ch := make(chan A, 1)
	go func() {
		ch <- f()
	}()
	return ch
}

func Async0(f func()) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		f()
		ch <- struct{}{}
	}()
	return ch
}

func Async1[A any](f func() A) <-chan A {
	return Async(f)
}

func Async2[A any, B any](f func() (A, B)) <-chan Tuple2[A, B] {
	ch := make(chan Tuple2[A, B], 1)
	go func() {
		ch <- T2(f())
	}()
	return ch
}
func Async3[A any, B any](ch chan<- Tuple2[A, B], f func() (A, B)) {
	go func() {
		ch <- T2(f())
	}()
}
