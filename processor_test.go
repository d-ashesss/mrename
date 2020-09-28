package main

import (
	"errors"
	"reflect"
	"testing"
)

type MemoryOutput map[string]string

func (m MemoryOutput) Put(name, result string) error {
	m[name] = result
	return nil
}

type MemoryFile struct {
	name string
	content string
}

func (m MemoryFile) Name() string {
	return m.name
}

type MapFileProvider map[string]string

func (m MapFileProvider) GetFiles() ([]FileInfo, error) {
	files := make([]FileInfo, 0, len(m))
	for name, _ := range m {
		files = append(files, MemoryFile{name: name})
	}
	return files, nil
}

type ErrorFileProvider struct {
	Error error
}

func (e ErrorFileProvider) GetFiles() ([]FileInfo, error) {
	return nil, e.Error
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

	t.Run("ProviderError", func(t *testing.T) {
		testError := errors.New("test")
		fileProvider := ErrorFileProvider{Error: testError}
		err := processor.Process(fileProvider)
		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}
	})
}
