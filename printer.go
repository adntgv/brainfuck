package main

import (
	"fmt"
)

type Printer interface {
	Printf(format string, values ...interface{})
}

func NewPrinter() Printer {
	return &defaultPrinter{}
}

type defaultPrinter struct {
}

func (p *defaultPrinter) Printf(format string, values ...interface{}) {
	fmt.Printf(format, values...)
}
