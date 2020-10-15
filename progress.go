package main

import (
	"log"
)

type ProgressAggregator interface {
	AddResult(name, result string)
}

type LoggedProgress struct {
	Logger  *log.Logger
	Verbose bool
}

func (o LoggedProgress) AddResult(name, result string) {
	if o.Verbose {
		o.Logger.Printf("%s %s\n", name, result)
	}
}
