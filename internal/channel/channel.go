package channel

func Forward[T any](src <-chan T, dest chan<- T) {
	for srcval := range src {
		dest <- srcval
	}
}
