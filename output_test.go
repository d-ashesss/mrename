package main

import "testing"

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
