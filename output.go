package main

import "encoding/json"

type Output interface {
	Format(progress ProgressAggregator) string
}

type JsonOutput struct {
}

func (o JsonOutput) Format(progress ProgressAggregator) string {
	j, _ := json.MarshalIndent(progress.GetResults(), "", "  ")
	return string(j)
}
