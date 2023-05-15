package file

import (
	"github.com/d-ashesss/mrename/producer"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Converter interface {
	Convert(Info) (string, error)
	SetNext(Converter)
}

type ReaderProducer interface {
	Produce(io.Reader) (string, error)
}

type converterChain struct {
	next Converter
}

func (c *converterChain) SetNext(n Converter) {
	c.next = n
}

func (c *converterChain) convertNext(i Info, newName string) (string, error) {
	if c.next != nil {
		return c.next.Convert(&namedInfo{Info: i, name: newName})
	}
	return newName, nil
}

type contentConverter struct {
	converterChain
	producer ReaderProducer
}

func (c *contentConverter) Convert(i Info) (string, error) {
	reader, err := fs.Open(i.Path())
	if err != nil {
		return "", err
	}
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(reader)
	newName, err := c.producer.Produce(reader)
	if err != nil {
		return "", err
	}
	if ext := filepath.Ext(i.Name()); ext != "" {
		newName += ext
	}
	return c.convertNext(i, newName)
}

func NewMD5Converter() Converter {
	return &contentConverter{producer: producer.MD5{}}
}

func NewSHA1Converter() Converter {
	return &contentConverter{producer: producer.SHA1{}}
}

type toLowerConverter struct {
	converterChain
}

func (c *toLowerConverter) Convert(i Info) (string, error) {
	newName := strings.ToLower(i.Name())
	return c.convertNext(i, newName)
}

func NewToLowerConverter() Converter {
	return &toLowerConverter{}
}

type jpeg2JpgConverter struct {
	converterChain
}

func (c *jpeg2JpgConverter) Convert(i Info) (string, error) {
	newName := i.Name()
	if regexp.MustCompile(`(?i)\.jpeg$`).MatchString(newName) {
		ext := path.Ext(i.Name())
		newName, _ = strings.CutSuffix(i.Name(), ext)
		newName += regexp.MustCompile(`(?i)e`).ReplaceAllString(ext, "")
	}
	return c.convertNext(i, newName)
}

func NewJpeg2JpgConverter() Converter {
	return &jpeg2JpgConverter{}
}
