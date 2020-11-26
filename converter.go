package main

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
)

type Converter interface {
	Convert(reader io.Reader) (string, error)
}

func convertByHash(reader io.Reader, h hash.Hash) (string, error) {
	_, err := io.Copy(h, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

type Md5Converter struct {
}

func (c Md5Converter) Convert(reader io.Reader) (string, error) {
	h := md5.New()
	return convertByHash(reader, h)
}
