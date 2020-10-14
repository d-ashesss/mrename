package main

import (
	"github.com/spf13/afero"
	"io"
	"path"
)

type FileInfo interface {
	Name() string
}

type FileProvider interface {
	GetFiles() ([]FileInfo, error)
	Open(info FileInfo) (io.Reader, error)
	Rename(info FileInfo, dstName string) error
}

type DirectoryFileProvider struct {
	Fs        afero.Fs
	Directory string
}

func (d DirectoryFileProvider) GetFiles() ([]FileInfo, error) {
	items, err := afero.ReadDir(d.Fs, d.Directory)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(items))
	for _, item := range items {
		if !item.IsDir() {
			files = append(files, item)
		}
	}
	return files, nil
}

func (d DirectoryFileProvider) Open(info FileInfo) (io.Reader, error) {
	filePath := path.Join(d.Directory, info.Name())
	file, err := d.Fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (d DirectoryFileProvider) Rename(info FileInfo, dstName string) error {
	filePath := path.Join(d.Directory, info.Name())
	dstPath := path.Join(d.Directory, dstName)
	return d.Fs.Rename(filePath, dstPath)
}
