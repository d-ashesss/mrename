package main

import (
	"fmt"
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var (
	dryRun bool
	verbose bool
)

func init() {
	flag.BoolVarP(&dryRun, "dry-run", "n", false, "Do not actually rename files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.Parse()
}

func main() {
	logger := log.New(os.Stderr, "", 0)
	progress := LoggedProgress{Logger: logger, Verbose: verbose}
	converter := Md5Converter{}
	fileProcessor := FileProcessor{Progress: progress, Converter: converter, DryRun: dryRun}
	processor := BulkProcessor{FileProcessor: &fileProcessor}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
