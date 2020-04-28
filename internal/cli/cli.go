package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/exports"
)

var skipImages *bool
var skipTimeseries *bool
var outputDirectory *string
var stdout *bool

//myUsage replaces default usage since it doesn't include information on non-flags
func myUsage() {
	fmt.Println("Extracts information from Innosat-MATS rac-files")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTIONS] rac-file ...\n", os.Args[0])
	fmt.Println()
	flag.PrintDefaults()
}

func getCallback(
	stdout bool,
	outputDirectory string,
	skipImages bool,
	skipTimeseries bool,
) (exports.Callback, exports.CallbackTeardown, error) {
	if outputDirectory == "" && !stdout {
		flag.Usage()
		fmt.Println("\nExpected an output directory")
		return nil, nil, errors.New("Invalid arguments")
	}
	if skipTimeseries && (skipImages || stdout) {
		fmt.Println("Nothing will be extracted, only validating integrity of rac-file(s)")
	}

	if stdout {
		callback, teardown := exports.StdoutCallbackFactory(os.Stdout, !skipTimeseries)
		return callback, teardown, nil
	}
	callback, teardown := exports.DiskCallbackFactory(
		outputDirectory,
		!skipImages,
		!skipTimeseries,
	)
	return callback, teardown, nil
}

func processFiles(
	extractor common.ExtractFunction,
	inputFiles []string,
	callback exports.Callback,
) error {
	batch := make([]common.StreamBatch, len(inputFiles))
	for n, filename := range inputFiles {
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}

		batch[n] = common.StreamBatch{
			Buf: f,
			Origin: common.OriginDescription{
				Name:           filename,
				ProcessingDate: time.Now(),
			},
		}

	}
	extractor(callback, batch...)
	return nil
}

func init() {
	skipImages = flag.Bool("skip-images", false, "Extract images from rac-files.\n(Default: false)")
	skipTimeseries = flag.Bool("skip-timeseries", false, "Extract timeseries from rac-files.\n(Default: false)")
	outputDirectory = flag.String("output", "", "Directory to place images and/or timeseries data")
	stdout = flag.Bool("stdout", false, "Output to standard out instead of to disk (only timeseries)\n(Default: false)")
	flag.Usage = myUsage
}

func main() {
	flag.Parse()
	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		flag.Usage()
		fmt.Println("\nNo rac-files supplied")
		os.Exit(1)
	}
	callback, teardown, err := getCallback(*stdout, *outputDirectory, *skipImages, *skipTimeseries)
	if err != nil {
		os.Exit(1)
	}
	err = processFiles(common.ExtractData, inputFiles, callback)
	if err != nil {
		log.Fatal(err)
	}
	teardown()
}
