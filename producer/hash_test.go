package producer_test

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/d-ashesss/mrename/producer"
	"strings"
	"sync"
	"testing"
)

type FailingReader struct {
	Error error
}

func (f FailingReader) Read(_ []byte) (int, error) {
	return 0, f.Error
}

func TestMD5_Produce(t *testing.T) {
	t.Run("produce", func(t *testing.T) {
		buf := strings.NewReader("testing content")
		p := producer.MD5{}
		got, err := p.Produce(buf)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "b91a4b2655c770f90410dc67dc407633"
		if expected != got {
			t.Errorf("Expected hash %q, got %q", expected, got)
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		strs := map[string]string{
			"first":   fmt.Sprintf("%x", md5.Sum([]byte("first"))),
			"second":  fmt.Sprintf("%x", md5.Sum([]byte("second"))),
			"third":   fmt.Sprintf("%x", md5.Sum([]byte("third"))),
			"fourth":  fmt.Sprintf("%x", md5.Sum([]byte("fourth"))),
			"fifth":   fmt.Sprintf("%x", md5.Sum([]byte("fifth"))),
			"sixth":   fmt.Sprintf("%x", md5.Sum([]byte("sixth"))),
			"seventh": fmt.Sprintf("%x", md5.Sum([]byte("seventh"))),
			"eights":  fmt.Sprintf("%x", md5.Sum([]byte("eights"))),
			"nineth":  fmt.Sprintf("%x", md5.Sum([]byte("nineth"))),
			"tenth":   fmt.Sprintf("%x", md5.Sum([]byte("tenth"))),
		}
		p := producer.MD5{}

		var wg sync.WaitGroup
		wg.Add(len(strs))

		for s, h := range strs {
			go func(s, expected string) {
				defer wg.Done()
				got, _ := p.Produce(strings.NewReader(s))
				if expected != got {
					t.Errorf("Expected %q hash to be %q, got %q", s, expected, got)
				}
			}(s, h)
		}
		wg.Wait()
	})

	t.Run("failing reader", func(t *testing.T) {
		testError := errors.New("fail")
		buffer := FailingReader{testError}
		p := producer.MD5{}
		_, err := p.Produce(buffer)
		if err == nil {
			t.Error("Expected error, none given")
		}
		if err != testError {
			t.Errorf("Expected test error, got %#v", err)
		}
	})
}
