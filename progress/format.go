package progress

import (
	"encoding/json"
	"fmt"
)

type ResultsAggregator interface {
	GetResults() map[string]string
}

func Format(format string, a ResultsAggregator) (string, error) {
	switch format {
	case "":
		return "", nil
	case "json":
		return FormatJSON(a)
	}
	return "", fmt.Errorf("invalid format %q", format)
}

func FormatJSON(a ResultsAggregator) (string, error) {
	j, err := json.MarshalIndent(a.GetResults(), "", "  ")
	if err != nil {
		return "", err
	}
	return string(j), nil
}
