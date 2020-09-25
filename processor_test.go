package main

import (
	"reflect"
	"testing"
)

func TestProcess (t *testing.T) {
	p := Processor{}
	got := p.Process([]string{"source1", "source2"})
	want := map[string]string{
		"source1": "source1",
		"source2": "source2",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func BenchmarkProcess(b *testing.B) {
	input := make([]string, 100)
	for i := 0; i < len(input); i++ {
		input[i] = "source"
	}
	p := Processor{}
	p.Process(input)
}
