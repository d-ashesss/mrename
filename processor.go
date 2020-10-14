package main

import "path"

type Processor struct {
	Output    ResultAggregator
	Converter Converter
	DryRun    bool
}

func (p *Processor) Process(provider FileProvider) error {
	files, err := provider.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		fp, _ := provider.Open(file)
		newName, _ := p.Converter.Convert(fp)
		if ext := path.Ext(file.Name()); ext != "" {
			newName += ext
		}
		var err error
		if !p.DryRun {
			err = provider.Rename(file, newName)
		}
		if err == nil {
			_ = p.Output.Put(file.Name(), newName)
		}
	}
	return nil
}
