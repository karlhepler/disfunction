package funk

import (
	"context"
	"fmt"
	"sync"

	"github.com/karlhepler/disfunction/internal/channel"
)

type ContextKey fmt.Stringer

func GetContextVal[ContextVal any](ctx context.Context, key ContextKey) (ContextVal, error) {
	val, ok := ctx.Value(key).(ContextVal)
	if !ok {
		return val, fmt.Errorf("invalid context key; key=%s", key)
	}
	return val, nil
}

func GoGetContextVal[ContextVal any](ctx context.Context, key ContextKey) (ContextVal, <-chan error) {
	errchan := make(chan error)
	outs, errs := channel.Async(func(outchan chan ContextVal, errchan chan error) {
		val, err := GetContextVal[ContextVal](ctx, key)
		if err != nil {
			errchan <- err
		}
		outchan <- val
	})

	var wg sync.WaitGroup
	channel.GoForward(ctx, &wg, errs, errchan)
	log := <-outs
	wg.Wait()

	return log, errchan
}
