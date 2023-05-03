package main

import (
	"testing"
)

type MemoryProgress map[string]string

func (m MemoryProgress) AddResult(name, result string) {
	m[name] = result
	return
}

func (m MemoryProgress) GetResults() map[string]string {
	return m
}

func TestJsonOutput_Format(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		progress := MemoryProgress{}
		output := JsonOutput{}
		json := output.Format(progress)
		expectedJson := "{}"
		if expectedJson != json {
			t.Error("Formatted JSON does not match expected JSON", json)
		}
	})

	t.Run("not empty input", func(t *testing.T) {
		progress := MemoryProgress{
			"first": "1st",
		}
		output := JsonOutput{}
		json := output.Format(progress)
		expectedJson := `{
  "first": "1st"
}`
		if expectedJson != json {
			t.Error("Formatted JSON does not match expected JSON", json, expectedJson)
		}
	})
}
