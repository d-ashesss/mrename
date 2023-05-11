package producer

import (
	"crypto/md5"
	"crypto/sha1"
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

func (p MD5) Produce(reader io.Reader) (string, error) {
	h := md5.New()
	return produceHash(reader, h)
}

type SHA1 struct {
}

func (p SHA1) Produce(reader io.Reader) (string, error) {
	h := sha1.New()
	return produceHash(reader, h)
}
