package main

import (
	"crypto/md5"
	"fmt"
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"os"
)

var (
	dryRun bool
)

func init() {
	flag.BoolVarP(&dryRun, "dry-run", "n", false, "Do not actually rename files.")
	flag.Parse()
}

func main() {
	output := TextOutput{os.Stdout}
	hash := md5.New()
	converter := HashConverter{Hash: hash}
	processor := Processor{Output: output, Converter: converter, DryRun: dryRun}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
