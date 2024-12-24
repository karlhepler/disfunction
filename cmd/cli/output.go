package main

import (
	"fmt"
	"log"

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
		log.Println("[ERROR] " + res.Message)
		return
	}

	fmt.Println(res.Message)
}
