package main

import (
	"os"
)

func main() {
	output := TextOutput{os.Stdout}
	processor := Processor{Output: output}
	processor.Process(os.Args[1:])
}
