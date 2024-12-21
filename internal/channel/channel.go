package channel

func Forward[T any](src <-chan T, dest chan<- T) {
	for val := range src {
		dest <- val
	}
}
