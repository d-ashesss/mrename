package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
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
}

func (m MemoryFile) Name() string {
	return m.name
}

type PlainConverter struct {
}

func (p PlainConverter) Convert(reader io.Reader) (string, error) {
	content, _ := ioutil.ReadAll(reader)
	return string(content), nil
}

type MapFileProvider map[string]string

func (m MapFileProvider) GetFiles() ([]FileInfo, error) {
	files := make([]FileInfo, 0, len(m))
	for name := range m {
		files = append(files, MemoryFile{name: name})
	}
	return files, nil
}

func (m MapFileProvider) Open(info FileInfo) (io.Reader, error) {
	content := m[info.Name()]
	return bytes.NewBufferString(content), nil
}

type ErrorFileProvider struct {
	GetError error
	OpenError error
}

func (e ErrorFileProvider) GetFiles() ([]FileInfo, error) {
	return nil, e.GetError
}

func (e ErrorFileProvider) Open(_ FileInfo) (io.Reader, error) {
	return nil, e.OpenError
}

func TestProcessor_Process(t *testing.T) {
	output := MemoryOutput{}
	converter := PlainConverter{}
	processor := Processor{Output: output, Converter: converter}
	fileProvider := MapFileProvider{
		"1st.txt": "first",
		"2nd.txt": "second",
		"3rd":     "third",
	}
	err := processor.Process(fileProvider)
	if err != nil {
		t.Errorf("Unexpected error %#v", err)
	}
	expected := MemoryOutput{
		"1st.txt": "first.txt",
		"2nd.txt": "second.txt",
		"3rd":     "third",
	}
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, got %v", expected, output)
	}

	t.Run("ProviderError", func(t *testing.T) {
		testError := errors.New("test")
		fileProvider := ErrorFileProvider{GetError: testError}
		err := processor.Process(fileProvider)
		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}
	})
}
