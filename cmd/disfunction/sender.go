package main

import (
	"fmt"

	"github.com/karlhepler/disfunction/pkg/handler"
)

type ConsoleSender struct{}

func NewConsoleSender() ConsoleSender {
	return ConsoleSender{}
}

func (console ConsoleSender) Send(res handler.DisfunctionRes) {
	fmt.Println(res.GoFunc.Line)
}
