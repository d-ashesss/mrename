package main

import (
	"bytes"
	"errors"
	"testing"
)

var bufferClosedError = errors.New("buffer closed")

type ClosedBuffer struct {
}

func (b ClosedBuffer) WriteString(_ string) (n int, err error) {
	return 0, bufferClosedError
}

func TestTextOutput_Put(t *testing.T) {
	buffer := bytes.NewBufferString("")
	output := TextOutput{buffer}
	_ = output.Put("name1", "result1")
	_ = output.Put("name2", "result2")
	expected := `name1 result1
name2 result2
`
	got := buffer.String()
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}

	t.Run("closed buffer", func(t *testing.T) {
		buffer := ClosedBuffer{}
		output := TextOutput{buffer}
		err := output.Put("name1", "result1")
		if err == nil {
			t.Errorf("Expected error")
		}
		if err != bufferClosedError {
			t.Errorf("Expected %#v, got %#v", bufferClosedError, err)
		}
	})
}
