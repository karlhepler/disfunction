package must

import (
	"os"
)

func Env(key string) (val string) {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic("[PANIC] required environment variable: " + key)
	}
	return val
}

func Val[T any](val T, err error) T {
	if err != nil {
		panic("[PANIC] required value: " + err.Error())
	}
	return val
}
