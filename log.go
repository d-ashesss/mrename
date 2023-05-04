package main

import (
	"github.com/d-ashesss/mrename/observer"
	"log"
)

type EventLogger struct {
	Verbose bool
}

func (l EventLogger) Notify(e observer.Event) {
	if l.Verbose && e.Name == "file.completed" {
		log.Printf("completed: %s: %s", e.File, e.Result)
	}
	if e.Name == "file.error" {
		log.Printf("error: %s: %s", e.File, e.Error)
	}
}
