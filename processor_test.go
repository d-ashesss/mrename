package main

import (
	"reflect"
	"testing"
)

type MemoryOutput map[string]string

func (m MemoryOutput) Put(name, result string) error {
	m[name] = result
	return nil
}

func TestProcessor_Process (t *testing.T) {
	output := MemoryOutput{}
	processor := Processor{Output: output}
	fileProvider := MapFileProvider{
		"1st.txt": "first",
		"2nd.txt": "second",
	}
	err := processor.Process(fileProvider)
	if err != nil {
		t.Errorf("Unexpected error %#v", err)
	}
	expected := MemoryOutput{
		"1st.txt": "1st.txt",
		"2nd.txt": "2nd.txt",
	}
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}
