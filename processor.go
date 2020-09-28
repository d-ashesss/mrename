package main

type Processor struct {
	Output ResultAggregator
}

func (p *Processor) Process(provider FileProvider) error {
	files, err := provider.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		_ = p.Output.Put(file.Name(), file.Name())
	}
	return nil
}
