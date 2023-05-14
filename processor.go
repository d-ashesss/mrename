package main

import (
	"github.com/d-ashesss/mrename/file"
	"github.com/d-ashesss/mrename/observer"
	"sync"
)

type Source interface {
	GetFiles() ([]file.Info, error)
}

type Target interface {
	Rename(info file.Info, newName string) error
}

type Processor struct {
	observer  *observer.Observer
	converter file.Converter
}

func NewProcessor(o *observer.Observer, c file.Converter) *Processor {
	return &Processor{
		observer:  o,
		converter: c,
	}
}

func (p *Processor) Process(source Source, target Target) error {
	files, err := source.GetFiles()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, f := range files {
		go func(f file.Info) {
			defer wg.Done()
			result, err := p.converter.Convert(f)
			if err != nil {
				p.observer.PublishError("file.error", f.Name(), err)
				return
			}
			if err := target.Rename(f, result); err != nil {
				p.observer.PublishError("file.error", f.Name(), err)
				return
			}
			p.observer.PublishResult("file.completed", f.Name(), result)
		}(f)
	}
	wg.Wait()
	return nil
}
