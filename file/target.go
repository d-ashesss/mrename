package file

import (
	"errors"
	"github.com/spf13/afero"
	"path/filepath"
)

var ErrNotDirectory = errors.New("target is not a directory")

// Target represents destination directory, where renamed files must be placed.
type Target struct {
	Path string
}

// CreateTarget instantiates new Target object, creating destination directory if it does not exist.
func CreateTarget(path string) (*Target, error) {
	if ok, err := afero.IsDir(fs, path); !ok && err == nil {
		return nil, ErrNotDirectory
	}
	if err := fs.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &Target{Path: path}, nil
}

// Rename moves files into the target directory with a new name.
func (t *Target) Rename(info Info, newName string) error {
	newPath := filepath.Join(t.Path, newName)
	return fs.Rename(info.Path(), newPath)
}

// VoidTarget does not perform any actions on files. Suitable for DryRun mode
type VoidTarget struct {
}

// Rename moves files into the target directory with a new name.
func (t *VoidTarget) Rename(_ Info, _ string) error {
	return nil
}
