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
