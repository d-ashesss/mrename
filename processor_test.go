package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"
)

type MemoryProgress map[string]string

func (m MemoryProgress) AddResult(name, result string) {
	m[name] = result
	return
}

func (m MemoryProgress) GetResults() map[string]string {
	return m
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

func (m MapFileProvider) MkDir(_ string) error {
	return nil
}

func (m MapFileProvider) Rename(info FileInfo, dstName string) error {
	content := m[info.Name()]
	delete(m, info.Name())
	m[dstName] = content
	return nil
}

type VolumeFileProvider int

func (v VolumeFileProvider) GetFiles() ([]FileInfo, error) {
	files := make([]FileInfo, 0, v)
	for i := 0; i < int(v); i++ {
		files = append(files, MemoryFile{name: "test"})
	}
	return files, nil
}

func (v VolumeFileProvider) Open(_ FileInfo) (io.Reader, error) {
	panic("method not available")
}

func (v VolumeFileProvider) MkDir(_ string) error {
	return nil
}

func (v VolumeFileProvider) Rename(_ FileInfo, _ string) error {
	panic("method not available")
}

type ErrorFileProvider struct {
	Files       MapFileProvider
	GetError    error
	OpenError   error
	MkDirError  error
	RenameError error
}

func (e ErrorFileProvider) GetFiles() ([]FileInfo, error) {
	if e.GetError == nil {
		return e.Files.GetFiles()
	}
	return nil, e.GetError
}

func (e ErrorFileProvider) Open(info FileInfo) (io.Reader, error) {
	if e.OpenError == nil {
		return e.Files.Open(info)
	}
	return nil, e.OpenError
}

func (e ErrorFileProvider) MkDir(_ string) error {
	return e.MkDirError
}

func (e ErrorFileProvider) Rename(info FileInfo, dstName string) error {
	if e.RenameError == nil {
		return e.Files.Rename(info, dstName)
	}
	return e.RenameError
}

type TrackingProcessor struct {
	processed map[string]string
	m sync.Mutex
}

func (t *TrackingProcessor) Process(info FileInfo, targetDir string, _ FileProvider) {
	t.m.Lock()
	defer t.m.Unlock()
	t.processed[ info.Name() ] = filepath.Join(targetDir, info.Name())
}

type TimedProcessor time.Duration

func (p TimedProcessor) Process(_ FileInfo, _ string, _ FileProvider) {
	time.Sleep(time.Duration(p))
}

type processorTest struct {
	Progress ProgressAggregator
	Converter Converter
	LogBuffer *bytes.Buffer
	Logger *log.Logger
	Processor *FileProcessor
	FileProvider MapFileProvider
}

func setUpFileProcessorTest() processorTest {
	progress := MemoryProgress{}
	converter := PlainConverter{}
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)
	processor := FileProcessor{Progress: progress, Converter: converter, Logger: logger}
	fileProvider := MapFileProvider{
		"1st.txt": "first",
		"2nd.txt": "second",
		"3rd":     "third",
	}

	return processorTest{
		Progress:     progress,
		Converter:    converter,
		LogBuffer:    &logBuffer,
		Logger:       logger,
		Processor:    &processor,
		FileProvider: fileProvider,
	}
}

func TestFileProcessor_Process(t *testing.T) {
	test := setUpFileProcessorTest()
	fileInfo := MemoryFile{name: "1st.txt"}
	test.Processor.Process(fileInfo, "", test.FileProvider)

	expectedProgress := MemoryProgress{
		"1st.txt": "first.txt",
	}
	if !reflect.DeepEqual(expectedProgress, test.Progress) {
		t.Errorf("Expected progress %v, got %v", expectedProgress, test.Progress)
	}
	if _, ok := test.FileProvider["1st.txt"]; ok {
		t.Error("Original file was not removed")
	}
	if _, ok := test.FileProvider["first.txt"]; !ok {
		t.Error("File was not renamed")
	}
	expectedLog := ""
	if expectedLog != test.LogBuffer.String() {
		t.Errorf("Expected empty log, got %v", test.LogBuffer.String())
	}

	t.Run("target dir", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := MemoryFile{name: "1st.txt"}
		test.Processor.Process(fileInfo, "target", test.FileProvider)

		if _, ok := test.FileProvider["target/first.txt"]; !ok {
			t.Error("File was not moved into target directory")
		}
	})

	t.Run("no file extension", func(t *testing.T) {
		processorTest := setUpFileProcessorTest()
		fileInfo := MemoryFile{name: "3rd"}
		processorTest.Processor.Process(fileInfo, "", processorTest.FileProvider)

		expectedProgress := MemoryProgress{
			"3rd": "third",
		}
		if !reflect.DeepEqual(expectedProgress, processorTest.Progress) {
			t.Errorf("Expected progress %v, got %v", expectedProgress, processorTest.Progress)
		}
	})

	t.Run("dry run", func(t *testing.T) {
		test := setUpFileProcessorTest()
		test.Processor.DryRun = true
		fileInfo := MemoryFile{name: "1st.txt"}
		test.Processor.Process(fileInfo, "", test.FileProvider)

		expectedProgress := MemoryProgress{
			"1st.txt": "first.txt",
		}
		if !reflect.DeepEqual(expectedProgress, test.Progress) {
			t.Errorf("Expected %v, got %v", expectedProgress, test.Progress)
		}
		if _, ok := test.FileProvider["1st.txt"]; !ok {
			t.Error("Original file was removed")
		}
		if _, ok := test.FileProvider["first.txt"]; ok {
			t.Error("File was renamed")
		}
		expectedLog := ""
		if expectedLog != test.LogBuffer.String() {
			t.Errorf("Expected empty log, got %v", test.LogBuffer.String())
		}
	})

	t.Run("file open error", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := MemoryFile{name: "1st.txt"}
		testError := errors.New("test file can't be opened")
		fileProvider := ErrorFileProvider{OpenError: testError}
		test.Processor.Process(fileInfo, "", fileProvider)

		expectedProgress := MemoryProgress{}
		if !reflect.DeepEqual(expectedProgress, test.Progress) {
			t.Errorf("Expected progress %v, got %v", expectedProgress, test.Progress)
		}
		expectedLog := "1st.txt: test file can't be opened\n"
		if expectedLog != test.LogBuffer.String() {
			t.Errorf("Expected logged file open error, got %v", test.LogBuffer.String())
		}
	})

	t.Run("file rename error", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := MemoryFile{name: "1st.txt"}
		testError := errors.New("test file can't be renamed")
		fileProvider := ErrorFileProvider{Files: MapFileProvider{"1st.txt": "first"}, RenameError: testError}
		test.Processor.Process(fileInfo, "", fileProvider)

		expectedProgress := MemoryProgress{}
		if !reflect.DeepEqual(expectedProgress, test.Progress) {
			t.Errorf("Expected progress %v, got %v", expectedProgress, test.Progress)
		}
		expectedLog := "1st.txt: test file can't be renamed\n"
		if expectedLog != test.LogBuffer.String() {
			t.Errorf("Expected logged file rename error, got %v", test.LogBuffer.String())
		}
	})
}

func TestBulkProcessor_Process(t *testing.T) {
	fileProcessor := TrackingProcessor{processed: map[string]string{}}
	processor := BulkProcessor{FileProcessor: &fileProcessor}

	fileProvider := ErrorFileProvider{
		Files: MapFileProvider{
			"1st.txt": "first",
			"2nd.txt": "second",
			"3rd":     "third",
		},
	}
	err := processor.Process(fileProvider)
	if err != nil {
		t.Errorf("Expected no error, got %#v", err)
	}

	expectedProcessed := map[string]string{"1st.txt": "1st.txt", "2nd.txt": "2nd.txt", "3rd": "3rd"}
	if !reflect.DeepEqual(expectedProcessed, fileProcessor.processed) {
		t.Errorf("Expected progress %v, got %v", expectedProcessed, fileProcessor.processed)
	}

	t.Run("target dir", func(t *testing.T) {
		fileProcessor := TrackingProcessor{processed: map[string]string{}}
		processor := BulkProcessor{FileProcessor: &fileProcessor, Target: "target"}

		fileProvider := ErrorFileProvider{
			Files: MapFileProvider{
				"1st.txt": "first",
				"2nd.txt": "second",
				"3rd":     "third",
			},
		}
		err := processor.Process(fileProvider)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}

		expectedProcessed := map[string]string{"1st.txt": "target/1st.txt", "2nd.txt": "target/2nd.txt", "3rd": "target/3rd"}
		if !reflect.DeepEqual(expectedProcessed, fileProcessor.processed) {
			t.Errorf("Expected progress %v, got %v", expectedProcessed, fileProcessor.processed)
		}
	})

	t.Run("make dir error", func(t *testing.T) {
		fileProcessor := TrackingProcessor{processed: map[string]string{}}
		processor := BulkProcessor{FileProcessor: &fileProcessor, Target: "target"}

		testError := errors.New("target dir can't be created")
		fileProvider := ErrorFileProvider{
			Files: MapFileProvider{
				"1st.txt": "first",
				"2nd.txt": "second",
				"3rd":     "third",
			},
			MkDirError: testError,
		}
		err := processor.Process(fileProvider)

		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}

		expectedProcessed := map[string]string{}
		if !reflect.DeepEqual(expectedProcessed, fileProcessor.processed) {
			t.Errorf("Expected progress %v, got %v", expectedProcessed, fileProcessor.processed)
		}
	})

	t.Run("get files error", func(t *testing.T) {
		fileProcessor := TrackingProcessor{processed: map[string]string{}}
		processor := BulkProcessor{FileProcessor: &fileProcessor}

		testError := errors.New("test files can't be listed")
		fileProvider := ErrorFileProvider{GetError: testError}
		err := processor.Process(fileProvider)

		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}

		expectedProcessed := map[string]string{}
		if !reflect.DeepEqual(expectedProcessed, fileProcessor.processed) {
			t.Errorf("Expected progress %v, got %v", expectedProcessed, fileProcessor.processed)
		}
	})
}

func BenchmarkBulkProcessor_Process(b *testing.B) {
	fileProcessor := TimedProcessor(100 * time.Microsecond)
	processor := BulkProcessor{FileProcessor: &fileProcessor}

	fileProvider := VolumeFileProvider(b.N)
	_ = processor.Process(fileProvider)
}
