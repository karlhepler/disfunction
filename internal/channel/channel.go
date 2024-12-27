package channel

import (
	"context"
	"sync"
)

func Forward[T any](ctx context.Context, src <-chan T, dest chan<- T) {
	ForEach(ctx, src, func(srcval T) {
		select {
		case dest <- srcval:
		case <-ctx.Done():
			return //context is canceled
		}
	})
}

func ForEach[T any](ctx context.Context, src <-chan T, callback func(T)) {
	for {
		select {
		case srcval, ok := <-src:
			if !ok {
				return // src channel is closed
			}
			callback(srcval)
		case <-ctx.Done():
			return // context is canceled
		}
	}
}

func GoForward[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Forward(ctx, src, dest)
	}()
}

func GoForEach[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, callback func(T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ForEach(ctx, src, callback)
	}()
}

// type Processor func[T any](context.Context, <-chan T) Async[T]
// ^ Go doesn't allow generic function types.
//	 I could probably define an interface instead.

func Async[T any](callback func(outchan chan T, errchan chan error)) (<-chan T, <-chan error) {
	outchan, errchan := make(chan T), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)
		callback(outchan, errchan)
	}()
	return outchan, errchan
}
