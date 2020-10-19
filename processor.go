package main

import (
	"log"
	"path"
)

type Processor interface {
	Process(info FileInfo, provider FileProvider)
}

type FileProcessor struct {
	Progress  ProgressAggregator
	Converter Converter
	Logger    *log.Logger
	DryRun    bool
}

func (f *FileProcessor) Process(info FileInfo, provider FileProvider)  {
	var err error
	file, err := provider.Open(info)
	if err != nil {
		f.Logger.Printf("%v: %v", info.Name(), err)
		return
	}
	newName, _ := f.Converter.Convert(file)
	if ext := path.Ext(info.Name()); ext != "" {
		newName += ext
	}
	if !f.DryRun {
		err = provider.Rename(info, newName)
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

func (p *BulkProcessor) Process(provider FileProvider) error {
	files, err := provider.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		p.FileProcessor.Process(file, provider)
	}
	return nil
}
