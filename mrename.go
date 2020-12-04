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
	target string
)

func init() {
	flag.BoolVarP(&dryRun, "dry-run", "n", false, "Do not actually rename files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.StringVarP(&target, "target", "t", "", "Specify the target directory")
	flag.Parse()
}

func main() {
	logger := log.New(os.Stderr, "", 0)
	progress := LoggedProgress{Logger: logger, Verbose: verbose}
	converter := Md5Converter{}
	fileProcessor := FileProcessor{Progress: progress, Converter: converter, Logger: logger, DryRun: dryRun}
	processor := BulkProcessor{FileProcessor: &fileProcessor, Target: target}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
