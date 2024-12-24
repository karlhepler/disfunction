package channel

import "sync"

func Forward[T any](src <-chan T, dest chan<- T) {
	ForEach(src, func(srcval T) {
		dest <- srcval
	})
}

func ForEach[T any](src <-chan T, cb func(T)) {
	for srcval := range src {
		cb(srcval)
	}
}

func GoForward[T any](wg *sync.WaitGroup, src <-chan T, dest chan<- T) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Forward(src, dest)
	}()
}

func GoForEach[T any](wg *sync.WaitGroup, src <-chan T, cb func(T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ForEach(src, cb)
	}()
}
