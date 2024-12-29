package channel

import (
	"context"
	"fmt"
	"sync"

	"github.com/karlhepler/disfunction/internal/funk"
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

type MapFunc[INPUT, OUTPUT any] func(context.Context, []INPUT, MapperFunc[INPUT, OUTPUT]) ([]OUTPUT, error)
type MapperFunc[INPUT, OUTPUT any] func(ctx context.Context, elem INPUT, index int, slice []INPUT) (OUTPUT, error)

func Map[INPUT, OUTPUT any](ctx context.Context, elems []INPUT, mapper MapperFunc[INPUT, OUTPUT]) ([]OUTPUT, error) {
	var errs error
	var outs []OUTPUT

	for i, elem := range elems {
		out, err := mapper(ctx, elem, i, elems)
		if err != nil {
			errs = fmt.Errorf("error mapping elem [%d]%v: %w", i, elem, err)
		}
		outs = append(outs, out)
	}
	return outs, errs
}

func GoFwd[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Fwd(ctx, src, dest)
	}()
}

func GoForEach[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, callback func(T, int, []T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TODO(karlhepler): ForEach needs to return an error
		// TODO(karlhepler): GoForEach needs to return an errchan
		ForEach(ctx, src, callback)
	}()
}

func GoMap[INPUT, OUTPUT any](ctx context.Context, wg *sync.WaitGroup, elems <-chan INPUT, mapper MapperFunc[INPUT, OUTPUT]) (<-chan OUTPUT, <-chan error) {
	outchan, errchan := make(chan OUTPUT), make(chan error)
	GoForEach[INPUT](ctx, wg, elems, func(elem INPUT, index i, slice []INPUT) {
		out, err := mapper(ctx, 
	})
	return outchan, errchan
	// return Async[OUTPUT](func(outchan chan OUTPUT, errchan chan error) {
	// }, AsyncWithWaitGroup(wg))
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
