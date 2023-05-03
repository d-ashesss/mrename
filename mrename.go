package main

import (
	"errors"
	"fmt"
	"github.com/d-ashesss/mrename/file"
	"github.com/d-ashesss/mrename/observer"
	flag "github.com/spf13/pflag"
	"log"
)

var (
	dryRun       bool
	verbose      bool
	targetDir    string
	outputFormat string
)

func init() {
	flag.BoolVarP(&dryRun, "dry-run", "n", false, "Do not actually rename files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.StringVarP(&targetDir, "target", "t", "", "Specify the target directory")
	flag.StringVarP(&outputFormat, "output-format", "o", "", "Output renaming results in specified format")
	flag.Parse()
}

func main() {
	log.SetFlags(0)
	progress := NewProgressAggregator()
	obsrvr := observer.New()
	obsrvr.AddSubscriber(progress)
	obsrvr.AddSubscriber(Log{Verbose: verbose})
	var converter Converter
	switch flag.Arg(0) {
	case "md5":
		converter = Md5Converter{}
		break
	default:
		log.Fatal("error: invalid converter: ", flag.Arg(0))
	}
	fileProcessor := FileProcessor{Observer: obsrvr, Converter: converter}
	processor := BulkProcessor{FileProcessor: &fileProcessor}
	source := file.NewSource(".")

	var target Target
	if dryRun {
		target = &file.VoidTarget{}
	} else {
		var err error
		target, err = file.CreateTarget(targetDir)
		if err != nil {
			log.Fatal("error: ", err)
		}
	}

	if err := processor.Process(source, target); err != nil {
		log.Fatal("error: ", err)
	}
	output, err := formatOutput(progress)
	if err != nil {
		log.Fatal("error: ", err)
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
