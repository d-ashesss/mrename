package main

import (
	"log"
	"sync"
)

type ProgressAggregator interface {
	AddResult(name, result string)
	GetResults() map[string]string
}

func NewLoggedProgress(logger *log.Logger, verbose bool) LoggedProgress {
	return LoggedProgress{
		Logger:  logger,
		Verbose: verbose,
		results: make(map[string]string),
		mutex: new(sync.Mutex),
	}
}

type LoggedProgress struct {
	Logger  *log.Logger
	Verbose bool
	results map[string]string
	mutex   *sync.Mutex
}

func (o LoggedProgress) AddResult(name, result string) {
	o.mutex.Lock()
	o.results[name] = result
	o.mutex.Unlock()
	if o.Verbose {
		o.Logger.Printf("%s %s\n", name, result)
	}
}

func (o LoggedProgress) GetResults() map[string]string {
	return o.results
}
