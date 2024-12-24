package cout

import (
	"fmt"
	"log"

	"github.com/karlhepler/disfunction/pkg/disfunction"
)

func Log(a string) {
	log.Println(a)
}

func Send(res disfunction.RandomMsg) {
	msg := res.Message
	if res.Status == disfunction.StatusError {
		msg = "[ERROR] " + msg
	}
	fmt.Println(msg)
}
