package main

import (
	"fmt"
	"log"
	"os"

	"github.com/karlhepler/disfunction/pkg/disfunction"
)

type Output struct {
	//
}

func (out Output) Log(a string) {
	log.Println(a)
}

func (out Output) Send(res disfunction.RandomMsg) {
	if res.Status == disfunction.StatusError {
		fmt.Fprintln(os.Stderr, "[ERROR] "+res.Message)
		return
	}

	fmt.Println(res.Message)
}
