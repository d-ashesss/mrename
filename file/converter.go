package file

import (
	"github.com/d-ashesss/mrename/producer"
	"io"
	"path/filepath"
)

type Converter interface {
	Convert(Info) (string, error)
}

type Producer interface {
	Produce(io.Reader) (string, error)
}

type contentConverter struct {
	producer Producer
}

func (c *contentConverter) Convert(i Info) (string, error) {
	reader, err := fs.Open(i.Path())
	if err != nil {
		return "", err
	}
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(reader)
	result, err := c.producer.Produce(reader)
	if err != nil {
		return "", err
	}
	if ext := filepath.Ext(i.Name()); ext != "" {
		result += ext
	}
	return result, nil
}

func NewMD5Converter() Converter {
	return &contentConverter{producer: producer.MD5{}}
}

func NewSHA1Converter() Converter {
	return &contentConverter{producer: producer.SHA1{}}
}
