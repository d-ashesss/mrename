package file

import (
	"errors"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
)

var ErrNotDirectory = errors.New("target is not a directory")

// Target represents destination directory, where files with new names must be placed.
type Target interface {
	// Acquire receives a file with a new name into target.
	Acquire(info Info, newName string) error
}

// MoveTarget renames original files into new ones.
type MoveTarget struct {
	Path string
}

// CreateMoveTarget instantiates new MoveTarget object, creating destination directory if it does not exist.
func CreateMoveTarget(path string) (Target, error) {
	if err := createTargetDir(path); err != nil {
		return nil, err
	}
	return &MoveTarget{Path: path}, nil
}

func (t *MoveTarget) Acquire(info Info, newName string) error {
	newPath := filepath.Join(t.Path, newName)
	return fs.Rename(info.Path(), newPath)
}

// CopyTarget creates a copy of original file with a new name.
type CopyTarget struct {
	Path string
}

// CreateCopyTarget instantiates new CopyTarget object, creating destination directory if it does not exist.
func CreateCopyTarget(path string) (Target, error) {
	if err := createTargetDir(path); err != nil {
		return nil, err
	}
	return &CopyTarget{Path: path}, nil
}

func (t *CopyTarget) Acquire(info Info, newName string) error {
	reader, err := fs.Open(info.Path())
	if err != nil {
		return err
	}
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(reader)
	newPath := filepath.Join(t.Path, newName)
	return afero.WriteReader(fs, newPath, reader)
}

// VoidTarget does not perform any actions on files. Suitable for DryRun mode
type VoidTarget struct {
}

func NewVoidTarget() Target {
	return &VoidTarget{}
}

func (t *VoidTarget) Acquire(_ Info, _ string) error {
	return nil
}

func createTargetDir(path string) error {
	if ok, err := afero.IsDir(fs, path); !ok && err == nil {
		return ErrNotDirectory
	}
	if path == "" {
		return nil
	}
	return fs.MkdirAll(path, 0755)
}
