package main

import (
	"bytes"
	"errors"
	"github.com/d-ashesss/mrename/file"
	"github.com/d-ashesss/mrename/mocks"
	"io"
	"log"
	"path"
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

type StringInfo string

func (f StringInfo) Name() string {
	return path.Base(string(f))
}

func (f StringInfo) Path() string {
	return string(f)
}

type PlainConverter struct {
}

func (p PlainConverter) Convert(reader io.Reader) (string, error) {
	content, _ := io.ReadAll(reader)
	return string(content), nil
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (b *ClosingBuffer) Close() error {
	return nil
}

type MapFileProvider map[string]string

func (m MapFileProvider) GetFiles() ([]file.Info, error) {
	files := make([]file.Info, 0, len(m))
	for name := range m {
		files = append(files, StringInfo(name))
	}
	return files, nil
}

func (m MapFileProvider) Open(info file.Info) (io.ReadCloser, error) {
	content := m[info.Name()]
	return &ClosingBuffer{bytes.NewBufferString(content)}, nil
}

func (m MapFileProvider) Rename(info file.Info, dstName string) error {
	content := m[info.Name()]
	delete(m, info.Name())
	m[dstName] = content
	return nil
}

type VolumeFileProvider int

func (v VolumeFileProvider) GetFiles() ([]file.Info, error) {
	files := make([]file.Info, 0, v)
	for i := 0; i < int(v); i++ {
		files = append(files, StringInfo("test"))
	}
	return files, nil
}

func (v VolumeFileProvider) Open(_ file.Info) (io.ReadCloser, error) {
	panic("method not available")
}

func (v VolumeFileProvider) Rename(_ file.Info, _ string) error {
	panic("method not available")
}

type ErrorFileProvider struct {
	Files     MapFileProvider
	GetError  error
	OpenError error
}

func (e ErrorFileProvider) GetFiles() ([]file.Info, error) {
	if e.GetError == nil {
		return e.Files.GetFiles()
	}
	return nil, e.GetError
}

func (e ErrorFileProvider) Open(info file.Info) (io.ReadCloser, error) {
	if e.OpenError == nil {
		return e.Files.Open(info)
	}
	return nil, e.OpenError
}

type TrackingProcessor struct {
	processed map[string]string
	m         sync.Mutex
}

func (t *TrackingProcessor) Process(info file.Info, _ Source, _ Target) {
	t.m.Lock()
	defer t.m.Unlock()
	t.processed[info.Name()] = info.Path()
}

type TimedProcessor time.Duration

func (p TimedProcessor) Process(_ file.Info, _ Source, _ Target) {
	time.Sleep(time.Duration(p))
}

type processorTest struct {
	Progress     ProgressAggregator
	Converter    Converter
	LogBuffer    *bytes.Buffer
	Logger       *log.Logger
	Processor    *FileProcessor
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
	t.Run("file", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := StringInfo("1st.txt")
		target := mocks.NewTarget(t)
		target.On("Rename", fileInfo, "first.txt").Return(nil)
		test.Processor.Process(fileInfo, test.FileProvider, target)

		expectedProgress := MemoryProgress{
			"1st.txt": "first.txt",
		}
		if !reflect.DeepEqual(expectedProgress, test.Progress) {
			t.Errorf("Expected progress %v, got %v", expectedProgress, test.Progress)
		}
		expectedLog := ""
		if expectedLog != test.LogBuffer.String() {
			t.Errorf("Expected empty log, got %v", test.LogBuffer.String())
		}
	})

	t.Run("no file extension", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := StringInfo("3rd")
		target := mocks.NewTarget(t)
		target.On("Rename", fileInfo, "third").Return(nil)
		test.Processor.Process(fileInfo, test.FileProvider, target)

		expectedProgress := MemoryProgress{
			"3rd": "third",
		}
		if !reflect.DeepEqual(expectedProgress, test.Progress) {
			t.Errorf("Expected progress %v, got %v", expectedProgress, test.Progress)
		}
	})

	t.Run("file open error", func(t *testing.T) {
		test := setUpFileProcessorTest()
		fileInfo := StringInfo("1st.txt")
		testError := errors.New("test file can't be opened")
		fileProvider := ErrorFileProvider{OpenError: testError}
		target := mocks.NewTarget(t)
		test.Processor.Process(fileInfo, fileProvider, target)

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
		fileInfo := StringInfo("1st.txt")
		testError := errors.New("test file can't be renamed")
		target := mocks.NewTarget(t)
		target.On("Rename", fileInfo, "first.txt").Return(testError)
		test.Processor.Process(fileInfo, test.FileProvider, target)

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
	t.Run("files", func(t *testing.T) {
		fileProcessor := TrackingProcessor{processed: map[string]string{}}
		processor := BulkProcessor{FileProcessor: &fileProcessor}

		fileProvider := MapFileProvider{
			"1st.txt": "first",
			"2nd.txt": "second",
			"3rd":     "third",
		}
		target := mocks.NewTarget(t)
		err := processor.Process(fileProvider, target)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}

		expectedProcessed := map[string]string{"1st.txt": "1st.txt", "2nd.txt": "2nd.txt", "3rd": "3rd"}
		if !reflect.DeepEqual(expectedProcessed, fileProcessor.processed) {
			t.Errorf("Expected progress %v, got %v", expectedProcessed, fileProcessor.processed)
		}
	})

	t.Run("get files error", func(t *testing.T) {
		fileProcessor := TrackingProcessor{processed: map[string]string{}}
		processor := BulkProcessor{FileProcessor: &fileProcessor}

		testError := errors.New("test files can't be listed")
		fileProvider := ErrorFileProvider{GetError: testError}
		target := mocks.NewTarget(t)
		err := processor.Process(fileProvider, target)

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
	target := mocks.NewTarget(b)
	_ = processor.Process(fileProvider, target)
}
