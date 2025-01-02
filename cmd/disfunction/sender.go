package main

import (
	"github.com/karlhepler/disfunction/pkg/handler"
)

type ConsoleSender struct{}

func NewConsoleSender() ConsoleSender {
	return ConsoleSender{}
}

func (console ConsoleSender) Send(data handler.DisfunctionRes) {
	// fmt.Println(data.Patch)
}
