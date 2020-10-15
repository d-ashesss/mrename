package main

import (
	"bytes"
	"log"
	"testing"
)

func TestLoggedProgress_AddResult(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	output := LoggedProgress{Logger: logger, Verbose: false}
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
		output := LoggedProgress{Logger: logger, Verbose: true}
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
