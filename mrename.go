package main

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var (
	dryRun       bool
	verbose      bool
	target       string
	outputFormat string
)

func init() {
	flag.BoolVarP(&dryRun, "dry-run", "n", false, "Do not actually rename files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.StringVarP(&target, "target", "t", "", "Specify the target directory")
	flag.StringVarP(&outputFormat, "output-format", "o", "", "Output renaming results in specified format")
	flag.Parse()
}

func main() {
	logger := log.New(os.Stderr, "", 0)
	progress := NewLoggedProgress(logger, verbose)
	var converter Converter
	switch flag.Arg(0) {
	case "md5":
		converter = Md5Converter{}
		break
	default:
		logger.Fatalln("Invalid converter", flag.Arg(0))
	}
	fileProcessor := FileProcessor{Progress: progress, Converter: converter, Logger: logger, DryRun: dryRun}
	processor := BulkProcessor{FileProcessor: &fileProcessor, Target: target}
	fileProvider := DirectoryFileProvider{Fs: afero.NewOsFs(), Directory: "."}
	err := processor.Process(fileProvider)
	if err != nil {
		logger.Fatalln(err)
	}
	output, err := formatOutput(progress)
	if err != nil {
		logger.Fatalln(err)
	}
	if len(output) > 0 {
		fmt.Print(output)
	}
}

func formatOutput(progress ProgressAggregator) (string, error) {
	var output Output
	switch outputFormat {
	case "":
		return "", nil
	case "json":
		output = JsonOutput{}
	default:
		return "", errors.New("test")
	}
	j := output.Format(progress)
	return j, nil
}
