package progress

import (
	"github.com/d-ashesss/mrename/observer"
	"sync"
)

func NewAggregator() *Aggregator {
	return &Aggregator{
		results: make(map[string]string),
	}
}

type Aggregator struct {
	results map[string]string
	mutex   sync.Mutex
}

func (o *Aggregator) Notify(e observer.Event) {
	if e.Name == "file.completed" {
		o.AddResult(e.File, e.Result)
	}
}

func (o *Aggregator) AddResult(name, result string) {
	o.mutex.Lock()
	o.results[name] = result
	o.mutex.Unlock()
}

func (o *Aggregator) GetResults() map[string]string {
	return o.results
}
