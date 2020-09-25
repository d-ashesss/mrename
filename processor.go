package main

type Processor struct {
}

func (p *Processor) Process(input []string) map[string]string {
	result := map[string]string{}
	for _, file := range input {
		result[file] = file
	}
	return result
}
