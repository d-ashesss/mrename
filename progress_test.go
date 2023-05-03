package main

import (
	"github.com/d-ashesss/mrename/observer"
	"reflect"
	"sync"
	"testing"
)

func TestProgress_AddResult(t *testing.T) {
	progress := NewProgressAggregator()
	progress.AddResult("name1", "result1")
	results := progress.GetResults()
	expectedResults := map[string]string{
		"name1": "result1",
	}
	if !reflect.DeepEqual(expectedResults, results) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

func TestProgress_Notify(t *testing.T) {
	progress := NewProgressAggregator()
	obsrvr := observer.New()
	obsrvr.AddSubscriber(progress)
	obsrvr.PublishResult("file.completed", "name1", "result1")
	results := progress.GetResults()
	expectedResults := map[string]string{
		"name1": "result1",
	}
	if !reflect.DeepEqual(expectedResults, results) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

func TestProgress_concurrency(t *testing.T) {
	progress := NewProgressAggregator()

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			progress.AddResult("name1", "result1")
			wg.Done()
		}()
	}

	wg.Wait()
}
