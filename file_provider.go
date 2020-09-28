package main

import "github.com/spf13/afero"

type FileInfo interface {
	Name() string
}

type FileProvider interface {
	GetFiles() ([]FileInfo, error)
}

type DirectoryFileProvider struct {
	Fs afero.Fs
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
