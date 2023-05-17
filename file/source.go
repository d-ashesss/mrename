package file

import (
	"github.com/spf13/afero"
)

type Source interface {
	GetFiles() ([]Info, error)
}

// Source represents a source directory with means to list and read files in it.
type fsSource struct {
	Path string
}

// NewSource instantiates a new Source.
func NewSource(path string) Source {
	return &fsSource{Path: path}
}

func (s *fsSource) GetFiles() ([]Info, error) {
	items, err := afero.ReadDir(fs, s.Path)
	if err != nil {
		return nil, err
	}
	files := make([]Info, 0, len(items))
	for _, item := range items {
		if !item.IsDir() {
			i := &osInfo{FileInfo: item, path: s.Path}
			files = append(files, i)
		}
	}
	return files, nil
}
