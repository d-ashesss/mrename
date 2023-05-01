package main

import (
	"github.com/d-ashesss/mrename/file"
	"io"
	"log"
	"path/filepath"
)

type Source interface {
	GetFiles() ([]file.Info, error)
	Open(i file.Info) (io.ReadCloser, error)
}

type Target interface {
	Rename(info file.Info, newName string) error
}

type Processor interface {
	Process(info file.Info, provider Source, target Target)
}

type FileProcessor struct {
	Progress  ProgressAggregator
	Converter Converter
	Logger    *log.Logger
	DryRun    bool
}

func (f *FileProcessor) Process(info file.Info, source Source, target Target) {
	var err error
	reader, err := source.Open(info)
	if err != nil {
		f.Logger.Printf("%v: %v", info.Name(), err)
		return
	}
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(reader)
	newName, _ := f.Converter.Convert(reader)
	if ext := filepath.Ext(info.Name()); ext != "" {
		newName += ext
	}
	if !f.DryRun {
		err = target.Rename(info, newName)
	}
	if err != nil {
		f.Logger.Printf("%v: %v", info.Name(), err)
		return
	}
	f.Progress.AddResult(info.Name(), newName)
}

type BulkProcessor struct {
	FileProcessor Processor
}

func (p *BulkProcessor) Process(source Source, target Target) error {
	files, err := source.GetFiles()
	if err != nil {
		return err
	}

	resultChannel := make(chan bool)
	for _, f := range files {
		go func(file file.Info) {
			p.FileProcessor.Process(file, source, target)
			resultChannel <- true
		}(f)
	}
	for i := 0; i < len(files); i++ {
		<-resultChannel
	}
	return nil
}
