package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/exports"
	"github.com/innosat-mats/rac-extract-payload/internal/extractors"
)

var skipImages *bool
var skipTimeseries *bool
var outputDirectory *string
var stdout *bool
var aws *bool

//myUsage replaces default usage since it doesn't include information on non-flags
func myUsage() {
	fmt.Println("Extracts information from Innosat-MATS rac-files")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTIONS] rac-file ...\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println(`
The tool can be used to scan rac files for contents. Use the -stdout flag and
use command line tools to scan for interesting information, e.g.:
	rac -stdout my.rac  | grep STAT

Tip for finding parsing errors:
	rac -stdout my.rac | grep -E -e".*Error:[^<:]+" -o

or if you want the Buffer contents which can be rather large if you are unlucky:
	rac -stdout my.rac | grep -E -e".*Error:[^<]+" -o
	`)
}

func getCallback(
	toStdout bool,
	toAws bool,
	outputDirectory string,
	skipImages bool,
	skipTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown, error) {
	if outputDirectory == "" && !toStdout && !toAws {
		flag.Usage()
		fmt.Println("\nExpected an output directory")
		return nil, nil, errors.New("Invalid arguments")
	}
	if skipTimeseries && (skipImages || toStdout) {
		fmt.Println("Nothing will be extracted, only validating integrity of rac-file(s)")
	}

	if toStdout {
		callback, teardown := exports.StdoutCallbackFactory(os.Stdout, !skipTimeseries)
		return callback, teardown, nil
	} else if toAws {
		callback, teardown := exports.AWSS3CallbackFactory(!skipImages, !skipTimeseries, wg)
		return callback, teardown, nil
	}
	callback, teardown := exports.DiskCallbackFactory(
		outputDirectory,
		!skipImages,
		!skipTimeseries,
		wg,
	)
	return callback, teardown, nil
}

func processFiles(
	extractor extractors.ExtractFunction,
	inputFiles []string,
	callback common.Callback,
) error {
	batch := make([]extractors.StreamBatch, len(inputFiles))
	for n, filename := range inputFiles {
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}
		batch[n] = extractors.StreamBatch{
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
	aws = flag.Bool("aws", false, "Output to aws instead of disk (requires credentials and permissions)")
	flag.Usage = myUsage
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		flag.Usage()
		log.Fatal("No rac-files supplied")
	}
	callback, teardown, err := getCallback(*stdout, *aws, *outputDirectory, *skipImages, *skipTimeseries, &wg)
	if err != nil {
		log.Fatal(err)
	}
	err = processFiles(extractors.ExtractData, inputFiles, callback)
	if err != nil {
		log.Fatal(err)
	}
	teardown()
}
