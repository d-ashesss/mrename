package main

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

func TestLoggedProgress_AddResult(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	output := NewLoggedProgress(logger, false)
	output.AddResult("name1", "result1")
	output.AddResult("name2", "result2")
	expected := ""
	got := buffer.String()
	if expected != got {
		t.Errorf("Expected empty log, got %v", got)
	}

	t.Run("verbose", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := log.New(&buffer, "", 0)
		output := NewLoggedProgress(logger, true)
		output.AddResult("name1", "result1")
		output.AddResult("name2", "result2")
		expected := `name1 result1
name2 result2
`
		got := buffer.String()
		if expected != got {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})
}

func TestLoggedProgress_GetResults(t *testing.T) {
	output := NewLoggedProgress(nil, false)
	output.AddResult("name1", "result1")
	results := output.GetResults()
	expectedResults := map[string]string {
		"name1": "result1",
	}
	if !reflect.DeepEqual(expectedResults, results) {
		t.Errorf("Expected results %v, got %v", expectedResults, results)
	}
}
