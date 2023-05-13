package main

import (
	"fmt"
	"github.com/d-ashesss/mrename/file"
	"github.com/d-ashesss/mrename/observer"
	"github.com/d-ashesss/mrename/progress"
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
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	progrss := progress.NewAggregator()
	obsrvr := observer.New()
	obsrvr.AddSubscriber(progrss)
	obsrvr.AddSubscriber(EventLogger{Verbose: verbose})

	var converter file.Converter
	switch flag.Arg(0) {
	case "md5":
		converter = file.NewMD5Converter()
	case "sha1":
		converter = file.NewSHA1Converter()
	case "tolower":
		converter = file.NewToLowerConverter()
	case "jpeg2jpg":
		converter = file.NewJpeg2JpgConverter()
	default:
		log.Fatalf("error: invalid converter: %q", flag.Arg(0))
	}
	processor := NewProcessor(obsrvr, converter)
	source := file.NewSource(".")

	var target Target
	if dryRun {
		target = &file.VoidTarget{}
	} else {
		var err error
		target, err = file.CreateTarget(targetDir)
		if err != nil {
			log.Fatal("error: create target: ", err)
		}
	}

	if err := processor.Process(source, target); err != nil {
		log.Fatal("error: process: ", err)
	}

	output, err := progress.Format(outputFormat, progrss)
	if err != nil {
		log.Fatal("error: output: ", err)
	}
	fmt.Print(output)
}
