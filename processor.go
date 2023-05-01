package main

import (
	"github.com/d-ashesss/mrename/file"
	"io"
	"log"
	"path/filepath"
)

type Target interface {
	Rename(info file.Info, newName string) error
}

type Processor interface {
	Process(info FileInfo, targetDir string, provider FileProvider, target Target)
}

type FileProcessor struct {
	Progress  ProgressAggregator
	Converter Converter
	Logger    *log.Logger
	DryRun    bool
}

func (f *FileProcessor) Process(info FileInfo, targetDir string, provider FileProvider, target Target) {
	var err error
	file, err := provider.Open(info)
	if err != nil {
		f.Logger.Printf("%v: %v", info.Name(), err)
		return
	}
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(file)
	newName, _ := f.Converter.Convert(file)
	if ext := filepath.Ext(info.Name()); ext != "" {
		newName += ext
	}
	newPath := filepath.Join(targetDir, newName)
	if !f.DryRun {
		err = provider.Rename(info, newPath)
	}
	if err != nil {
		f.Logger.Printf("%v: %v", info.Name(), err)
		return
	}
	f.Progress.AddResult(info.Name(), newName)
}

type BulkProcessor struct {
	FileProcessor Processor
	Target        string
}

func (p *BulkProcessor) Process(provider FileProvider, target Target) error {
	files, err := provider.GetFiles()
	if err != nil {
		return err
	}

	resultChannel := make(chan bool)
	for _, file := range files {
		go func(file FileInfo) {
			p.FileProcessor.Process(file, p.Target, provider, target)
			resultChannel <- true
		}(file)
	}
	for i := 0; i < len(files); i++ {
		<-resultChannel
	}
	return nil
}
