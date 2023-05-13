package file

import (
	"errors"
	"github.com/d-ashesss/mrename/producer"
	"io"
	"path"
	"path/filepath"
	"strings"
)

var ErrFileSkipped = errors.New("file skipped")

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

type toLowerConverter struct {
}

func (c *toLowerConverter) Convert(i Info) (string, error) {
	return strings.ToLower(i.Name()), nil
}

func NewToLowerConverter() Converter {
	return &toLowerConverter{}
}

type jpeg2JpgConverter struct {
}

func (c *jpeg2JpgConverter) Convert(i Info) (string, error) {
	ext := path.Ext(i.Name())
	if ext != ".jpeg" {
		return i.Name(), ErrFileSkipped
	}
	name, _ := strings.CutSuffix(i.Name(), ext)
	return name + ".jpg", nil
}

func NewJpeg2JpgConverter() Converter {
	return &jpeg2JpgConverter{}
}
