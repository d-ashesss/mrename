package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

type MemoryProgress map[string]string

func (m MemoryProgress) AddResult(name, result string) {
	m[name] = result
	return
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

func (m MapFileProvider) Rename(info FileInfo, dstName string) error {
	content := m[info.Name()]
	delete(m, info.Name())
	m[dstName] = content
	return nil
}

type ErrorFileProvider struct {
	GetError    error
	OpenError   error
	RenameError error
}

func (e ErrorFileProvider) GetFiles() ([]FileInfo, error) {
	return nil, e.GetError
}

func (e ErrorFileProvider) Open(_ FileInfo) (io.Reader, error) {
	return nil, e.OpenError
}

func (e ErrorFileProvider) Rename(_ FileInfo, _ string) error {
	return e.RenameError
}

func TestProcessor_Process(t *testing.T) {
	output := MemoryProgress{}
	converter := PlainConverter{}
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)
	processor := Processor{Progress: output, Converter: converter, Logger: logger}
	fileProvider := MapFileProvider{
		"1st.txt": "first",
		"2nd.txt": "second",
		"3rd":     "third",
	}
	err := processor.Process(fileProvider)
	if err != nil {
		t.Errorf("Unexpected processing error %#v", err)
	}
	expectedProgress := MemoryProgress{
		"1st.txt": "first.txt",
		"2nd.txt": "second.txt",
		"3rd":     "third",
	}
	if !reflect.DeepEqual(expectedProgress, output) {
		t.Errorf("Expected %v, got %v", expectedProgress, output)
	}
	if _, ok := fileProvider["1st.txt"]; ok {
		t.Error("Original file was not removed")
	}
	if _, ok := fileProvider["first.txt"]; !ok {
		t.Error("File was not renamed")
	}
	expectedLog := ""
	if expectedLog != logBuffer.String() {
		t.Errorf("Expected empty log, got %v", logBuffer.String())
	}

	t.Run("DryRun", func(t *testing.T) {
		output := MemoryProgress{}
		converter := PlainConverter{}
		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)
		processor := Processor{Progress: output, Converter: converter, DryRun: true, Logger: logger}
		fileProvider := MapFileProvider{
			"1st.txt": "first",
			"2nd.txt": "second",
			"3rd":     "third",
		}
		err := processor.Process(fileProvider)
		if err != nil {
			t.Errorf("Unexpected processing error %#v", err)
		}
		expectedProgress := MemoryProgress{
			"1st.txt": "first.txt",
			"2nd.txt": "second.txt",
			"3rd":     "third",
		}
		if !reflect.DeepEqual(expectedProgress, output) {
			t.Errorf("Expected %v, got %v", expectedProgress, output)
		}
		if _, ok := fileProvider["1st.txt"]; !ok {
			t.Error("Original file was removed")
		}
		if _, ok := fileProvider["first.txt"]; ok {
			t.Error("File was renamed")
		}
		expectedLog := ""
		if expectedLog != logBuffer.String() {
			t.Errorf("Expected empty log, got %v", logBuffer.String())
		}
	})

	t.Run("ProviderError", func(t *testing.T) {
		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)
		processor := Processor{Progress: output, Converter: converter, Logger: logger}
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
