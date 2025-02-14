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

func FwdToOutchan[T any](ctx context.Context, src <-chan T, outchan <-chan T) <-chan T {
	out2chan, _ := Async(func(out2chan chan T, _ chan error) {
		Fwd(ctx, src, out2chan)
		Fwd(ctx, outchan, out2chan)
	})
	return out2chan
}

func Filter[T any](ctx context.Context, in <-chan T, predicate func(T) bool) <-chan T {
	outs := make(chan T)
	go func() {
		defer close(outs)
		ForEach(ctx, in, func(val T) {
			if predicate(val) {
				outs <- val
			}
		})
	}()
	return outs
}

type MapperFunc[IN, OUT any] func(IN, chan<- OUT, chan<- error)

func Map[IN, OUT any](ctx context.Context, in <-chan IN, mapper MapperFunc[IN, OUT]) (<-chan OUT, <-chan error) {
	outs, errs := make(chan OUT), make(chan error)
	go func() {
		defer close(outs)
		defer close(errs)
		ForEach(ctx, in, func(val IN) {
			mapper(val, outs, errs)
		})
	}()
	return outs, errs
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

func Collect[T any](ctx context.Context, src <-chan T) []T {
	var out []T
	ForEach(ctx, src, func(val T) {
		out = append(out, val)
	})
	return out
}

func GoFwd[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Fwd(ctx, src, dest)
	}()
}

func GoFwdToOutchan[T any](ctx context.Context, wg *sync.WaitGroup, src <-chan T, outchan <-chan T) <-chan T {
	wg.Add(1)
	var out2chan <-chan T
	go func() {
		defer wg.Done()
		out2chan = FwdToOutchan(ctx, src, outchan)
	}()
	return out2chan
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

func NewOutchanOverride[T any](outchan <-chan T, override func(out T) (T, error)) (<-chan T, <-chan error) {
	return Async(func(out2chan chan T, errchan chan error) {
		for out := range outchan {
			out2, err := override(out)
			if err != nil {
				errchan <- err
			}
			out2chan <- out2
		}
	})
}

func SendOnOutchan[T any](data T, outchan <-chan T) <-chan T {
	out2chan, _ := Async(func(out2chan chan T, _ chan error) {
		out2chan <- data
		for out := range outchan {
			out2chan <- out
		}
	})
	return out2chan
}

func SendEachOnChannel[T any](data []T) <-chan T {
	datumchan := make(chan T)
	go func() {
		defer close(datumchan)
		for _, datum := range data {
			datumchan <- datum
		}
	}()
	return datumchan
}
