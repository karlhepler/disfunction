package main

import (
	"fmt"
)

type ConsoleSender struct{}

func NewConsoleSender() ConsoleSender {
	return ConsoleSender{}
}

func (console ConsoleSender) Send(res any) {
	fmt.Println(res)
}
