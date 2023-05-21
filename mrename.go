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
	copyFiles    bool
	dryRun       bool
	verbose      bool
	targetDir    string
	outputFormat string
)

func init() {
	flag.BoolVarP(&copyFiles, "copy", "c", false, "Copy files instead of renaming them")
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
	for i := flag.NArg() - 1; i >= 0; i-- {
		var c file.Converter
		switch flag.Arg(i) {
		case "md5":
			c = file.NewMD5Converter()
		case "sha1":
			c = file.NewSHA1Converter()
		case "tolower":
			c = file.NewToLowerConverter()
		case "toupper":
			c = file.NewToUpperConverter()
		case "jpeg2jpg":
			c = file.NewJpeg2JpgConverter()
		default:
			log.Fatalf("error: invalid converter: %q", flag.Arg(i))
		}
		c.SetNext(converter)
		converter = c
	}
	processor := NewProcessor(obsrvr, converter)
	source := file.NewSource(".")

	target, err := createTarget(targetDir)
	if err != nil {
		log.Fatal("error: create target: ", err)
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

func createTarget(path string) (file.Target, error) {
	if dryRun {
		return file.NewVoidTarget(), nil
	}
	if copyFiles {
		return file.CreateCopyTarget(path)
	}
	return file.CreateMoveTarget(path)
}
