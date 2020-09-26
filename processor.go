package main

type Processor struct {
	Output ResultAggregator
}

func (p *Processor) Process(input []string) {
	for _, file := range input {
		_ = p.Output.Put(file, file)
	}
}
