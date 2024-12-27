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

func ForEach[T any](ctx context.Context, src <-chan T, cb func(T)) {
	for {
		select {
		case srcval, ok := <-src:
			if !ok {
				return // src channel is closed
			}
			cb(srcval)
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

func GoForEach[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, cb func(T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ForEach(ctx, src, cb)
	}()
}
