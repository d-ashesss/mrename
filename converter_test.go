package main

import (
	"bytes"
	"errors"
	"testing"
)

type FailingReader struct {
	Error error
}

func (f FailingReader) Read(_ []byte) (int, error) {
	return 0, f.Error
}

func TestMd5Converter_Convert(t *testing.T) {
	buffer := bytes.NewBufferString("testing content")
	converter := Md5Converter{}
	got, _ := converter.Convert(buffer)
	expected := "b91a4b2655c770f90410dc67dc407633"
	if expected != got {
		t.Errorf("Expected hash %v, got %v", expected, got)
	}

	t.Run("reuse converter", func(t *testing.T) {
		buffer := bytes.NewBufferString("tstng cntnt")
		got, _ := converter.Convert(buffer)
		expected := "330a70b3938d00a605aaf18b44f4184f"
		if expected != got {
			t.Errorf("Expected hash %v, got %v", expected, got)
		}
	})

	t.Run("failing reader", func(t *testing.T) {
		testError := errors.New("fail")
		buffer := FailingReader{testError}
		_, err := converter.Convert(buffer)
		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}
	})
}
