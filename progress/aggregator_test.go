package progress_test

import (
	"github.com/d-ashesss/mrename/observer"
	"github.com/d-ashesss/mrename/progress"
	"reflect"
	"sync"
	"testing"
)

func TestAggregator_AddResult(t *testing.T) {
	aggregator := progress.NewAggregator()
	aggregator.AddResult("name1", "result1")
	results := aggregator.GetResults()
	expectedResults := map[string]string{
		"name1": "result1",
	}
	if !reflect.DeepEqual(expectedResults, results) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

func TestAggregator_Notify(t *testing.T) {
	aggregator := progress.NewAggregator()
	obsrvr := observer.New()
	obsrvr.AddSubscriber(aggregator)
	obsrvr.PublishResult("file.completed", "name1", "result1")
	results := aggregator.GetResults()
	expectedResults := map[string]string{
		"name1": "result1",
	}
	if !reflect.DeepEqual(expectedResults, results) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}

func TestAggregator_concurrency(t *testing.T) {
	aggregator := progress.NewAggregator()

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			aggregator.AddResult("name1", "result1")
			wg.Done()
		}()
	}

	wg.Wait()
}
