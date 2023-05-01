package file

import (
	"github.com/spf13/afero"
	"io"
)

// Source represents a source directory with means to list and read files in it.
type Source struct {
	Path string
}

// NewSource instantiates a new Source.
func NewSource(path string) *Source {
	return &Source{Path: path}
}

func (s *Source) GetFiles() ([]Info, error) {
	items, err := afero.ReadDir(fs, s.Path)
	if err != nil {
		return nil, err
	}
	files := make([]Info, 0, len(items))
	for _, item := range items {
		if !item.IsDir() {
			i := &info{FileInfo: item, path: s.Path}
			files = append(files, i)
		}
	}
	return files, nil
}

func (s *Source) Open(i Info) (io.ReadCloser, error) {
	return fs.Open(i.Path())
}
