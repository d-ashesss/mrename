package producer

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
)

func produceHash(reader io.Reader, h hash.Hash) (string, error) {
	_, err := io.Copy(h, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

type MD5 struct {
}

func (c MD5) Produce(reader io.Reader) (string, error) {
	h := md5.New()
	return produceHash(reader, h)
}
