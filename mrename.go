package main

import (
	"crypto/md5"
	"fmt"
	"github.com/spf13/afero"
	"os"
)

func main() {
	output := TextOutput{os.Stdout}
	hash := md5.New()
	converter := HashConverter{Hash: hash}
	processor := Processor{Output: output, Converter: converter}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
