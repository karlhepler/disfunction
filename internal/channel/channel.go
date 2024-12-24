package channel

import (
	"sync"
)

func Forward[T any](src <-chan T, dest chan<- T) {
	for srcval := range src {
		dest <- srcval
	}
}

func GoForward[T any](wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Forward(src, dest)
	}()
}
