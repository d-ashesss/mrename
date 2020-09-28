package main

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
)

func main() {
	output := TextOutput{os.Stdout}
	processor := Processor{Output: output}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
