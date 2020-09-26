package main

import (
	"fmt"
	"io"
)

type ResultAggregator interface {
	Put(name, result string) error
}

type TextOutput struct {
	Writer io.StringWriter
}

func (o TextOutput) Put(name, result string) error {
	_, err := o.Writer.WriteString(fmt.Sprintf("%s %s\n", name, result))
	return err
}
