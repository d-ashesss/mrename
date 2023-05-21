package file

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

type MD5Producer struct {
}

func (p MD5Producer) Produce(reader io.Reader) (string, error) {
	h := md5.New()
	return produceHash(reader, h)
}

type SHA1Producer struct {
}

func (p SHA1Producer) Produce(reader io.Reader) (string, error) {
	h := sha1.New()
	return produceHash(reader, h)
}
