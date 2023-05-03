package main

import (
	"errors"
	"github.com/d-ashesss/mrename/file"
	"github.com/d-ashesss/mrename/mocks"
	"github.com/d-ashesss/mrename/observer"
	"path"
	"testing"
)

type StringInfo string

func (f StringInfo) Name() string {
	return path.Base(string(f))
}

func (f StringInfo) Path() string {
	return string(f)
}

func TestProcessor_Process(t *testing.T) {
	t.Run("files", func(t *testing.T) {
		obsrvr := observer.New()
		subscriber := mocks.NewSubscriber(t)
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "1st.txt", Result: "fst"})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "2nd.txt", Result: "snd"})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "3rd", Result: "trd"})
		obsrvr.AddSubscriber(subscriber)
		converter := mocks.NewConverter(t)
		converter.On("Convert", StringInfo("1st.txt")).Return("fst", nil).Once()
		converter.On("Convert", StringInfo("2nd.txt")).Return("snd", nil).Once()
		converter.On("Convert", StringInfo("3rd")).Return("trd", nil).Once()
		processor := NewProcessor(obsrvr, converter)

		source := mocks.NewSource(t)
		source.On("GetFiles").Return([]file.Info{
			StringInfo("1st.txt"),
			StringInfo("2nd.txt"),
			StringInfo("3rd"),
		}, nil)
		target := mocks.NewTarget(t)
		target.On("Rename", StringInfo("1st.txt"), "fst").Return(nil).Once()
		target.On("Rename", StringInfo("2nd.txt"), "snd").Return(nil).Once()
		target.On("Rename", StringInfo("3rd"), "trd").Return(nil).Once()
		err := processor.Process(source, target)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
	})

	t.Run("file converter error", func(t *testing.T) {
		testError := errors.New("test file can't be opened")
		obsrvr := observer.New()
		subscriber := mocks.NewSubscriber(t)
		subscriber.On("Notify", observer.Event{Name: "file.error", File: "1st.txt", Error: testError})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "2nd.txt", Result: "snd"})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "3rd", Result: "trd"})
		obsrvr.AddSubscriber(subscriber)
		converter := mocks.NewConverter(t)
		converter.On("Convert", StringInfo("1st.txt")).Return("", testError).Once()
		converter.On("Convert", StringInfo("2nd.txt")).Return("snd", nil).Once()
		converter.On("Convert", StringInfo("3rd")).Return("trd", nil).Once()
		processor := NewProcessor(obsrvr, converter)

		source := mocks.NewSource(t)
		source.On("GetFiles").Return([]file.Info{
			StringInfo("1st.txt"),
			StringInfo("2nd.txt"),
			StringInfo("3rd"),
		}, nil)
		target := mocks.NewTarget(t)
		target.On("Rename", StringInfo("2nd.txt"), "snd").Return(nil).Once()
		target.On("Rename", StringInfo("3rd"), "trd").Return(nil).Once()
		err := processor.Process(source, target)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
	})

	t.Run("file rename error", func(t *testing.T) {
		testError := errors.New("test file can't be renamed")
		obsrvr := observer.New()
		subscriber := mocks.NewSubscriber(t)
		subscriber.On("Notify", observer.Event{Name: "file.error", File: "1st.txt", Error: testError})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "2nd.txt", Result: "snd"})
		subscriber.On("Notify", observer.Event{Name: "file.completed", File: "3rd", Result: "trd"})
		obsrvr.AddSubscriber(subscriber)
		converter := mocks.NewConverter(t)
		converter.On("Convert", StringInfo("1st.txt")).Return("fst", nil).Once()
		converter.On("Convert", StringInfo("2nd.txt")).Return("snd", nil).Once()
		converter.On("Convert", StringInfo("3rd")).Return("trd", nil).Once()
		processor := NewProcessor(obsrvr, converter)

		source := mocks.NewSource(t)
		source.On("GetFiles").Return([]file.Info{
			StringInfo("1st.txt"),
			StringInfo("2nd.txt"),
			StringInfo("3rd"),
		}, nil)
		target := mocks.NewTarget(t)
		target.On("Rename", StringInfo("1st.txt"), "fst").Return(testError).Once()
		target.On("Rename", StringInfo("2nd.txt"), "snd").Return(nil).Once()
		target.On("Rename", StringInfo("3rd"), "trd").Return(nil).Once()
		err := processor.Process(source, target)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
	})

	t.Run("get files error", func(t *testing.T) {
		processor := Processor{}

		testError := errors.New("test files can't be listed")
		source := mocks.NewSource(t)
		source.On("GetFiles").Return(nil, testError)
		target := mocks.NewTarget(t)
		err := processor.Process(source, target)

		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}
	})
}
