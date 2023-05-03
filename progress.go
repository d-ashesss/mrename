package main

import (
	"github.com/d-ashesss/mrename/observer"
	"sync"
)

func NewProgressAggregator() *Progress {
	return &Progress{
		results: make(map[string]string),
	}
}

type Progress struct {
	results map[string]string
	mutex   sync.Mutex
}

func (o *Progress) Notify(e observer.Event) {
	if e.Name == "file.completed" {
		o.AddResult(e.File, e.Result)
	}
}

func (o *Progress) AddResult(name, result string) {
	o.mutex.Lock()
	o.results[name] = result
	o.mutex.Unlock()
}

func (o *Progress) GetResults() map[string]string {
	return o.results
}
