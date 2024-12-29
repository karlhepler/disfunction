package channel

import (
	"context"
	"sync"
)

func Fwd[T any](ctx context.Context, src <-chan T, dest chan<- T) {
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

func GoFwd[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Fwd(ctx, src, dest)
	}()
}

func GoForEach[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, callback func(T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ForEach(ctx, src, callback)
	}()
}

type AsyncFunc[T any] func(callback func(outchan chan T, errchan chan error) (<-chan T, <-chan error))
type ProcessorFunc[T any] func(context.Context, <-chan T) AsyncFunc[T]

func Async[T any](callback func(outchan chan T, errchan chan error)) (<-chan T, <-chan error) {
	outchan, errchan := make(chan T), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)
		callback(outchan, errchan)
	}()
	return outchan, errchan
}
