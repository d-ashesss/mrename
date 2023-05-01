package file

import (
	"os"
	"path"
)

type Info interface {
	Name() string
	Path() string
}

type info struct {
	os.FileInfo
	path string
}

func (i *info) Path() string {
	return path.Join(i.path, i.Name())
}
