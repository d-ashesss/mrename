package main

import "encoding/json"

type ProgressAggregator interface {
	AddResult(name, result string)
	GetResults() map[string]string
}

type Output interface {
	Format(progress ProgressAggregator) string
}

type JsonOutput struct {
}

func (o JsonOutput) Format(progress ProgressAggregator) string {
	j, _ := json.MarshalIndent(progress.GetResults(), "", "  ")
	return string(j)
}
