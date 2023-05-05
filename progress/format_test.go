package progress_test

import (
	"github.com/d-ashesss/mrename/progress"
	"testing"
)

type MapResults map[string]string

func (m MapResults) GetResults() map[string]string {
	return m
}

func TestJsonOutput_Format(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		results := MapResults{}
		json, err := progress.FormatJSON(results)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		expectedJson := "{}"
		if expectedJson != json {
			t.Error("Formatted JSON does not match expected JSON", json)
		}
	})

	t.Run("not empty input", func(t *testing.T) {
		results := MapResults{
			"first": "1st",
		}
		json, err := progress.FormatJSON(results)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		expectedJson := `{
  "first": "1st"
}`
		if expectedJson != json {
			t.Error("Formatted JSON does not match expected JSON", json, expectedJson)
		}
	})
}
