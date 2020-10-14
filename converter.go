package main

import (
	"fmt"
	"hash"
	"io"
)

type Converter interface {
	Convert(reader io.Reader) (string, error)
}

type HashConverter struct {
	Hash hash.Hash
}

func (h HashConverter) Convert(reader io.Reader) (string, error) {
	digest := h.Hash
	digest.Reset()
	_, err := io.Copy(digest, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}
