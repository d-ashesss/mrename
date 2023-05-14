package file

import (
	"os"
	"path"
)

type Info interface {
	Name() string
	Path() string
}

type osInfo struct {
	os.FileInfo
	path string
}

func (i *osInfo) Path() string {
	return path.Join(i.path, i.Name())
}

type namedInfo struct {
	Info
	name string
}

func (i *namedInfo) Name() string {
	return i.name
}
