package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/awstools"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/exports"
	"github.com/innosat-mats/rac-extract-payload/internal/extractors"
)

var skipImages *bool
var skipTimeseries *bool
var project *string
var stdout *bool
var aws *bool
var awsDescription *string

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
	project string,
	skipImages bool,
	skipTimeseries bool,
	awsDescription string,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown, error) {
	if project == "" && !toStdout {
		flag.Usage()
		fmt.Println("\nExpected a project")
		return nil, nil, errors.New("Invalid arguments")
	}
	if skipTimeseries && (skipImages || toStdout) {
		fmt.Println("Nothing will be extracted, only validating integrity of rac-file(s)")
	}

	if toStdout {
		callback, teardown := exports.StdoutCallbackFactory(os.Stdout, !skipTimeseries)
		return callback, teardown, nil
	} else if toAws {
		callback, teardown := exports.AWSS3CallbackFactory(
			awstools.AWSUpload,
			project,
			awsDescription,
			!skipImages,
			!skipTimeseries,
			wg,
		)
		return callback, teardown, nil
	}
	callback, teardown := exports.DiskCallbackFactory(
		project,
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
	project = flag.String("project", "", "Name for experiments, when outputting to disk a directory will be created with this name, when sending to AWS files will have this as a prefix")
	stdout = flag.Bool("stdout", false, "Output to standard out instead of to disk (only timeseries)\n(Default: false)")
	aws = flag.Bool("aws", false, "Output to aws instead of disk (requires credentials and permissions)")
	awsDescription = flag.String("description", "", "Path to a file containing a project description to be uploaded to AWS")
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
	callback, teardown, err := getCallback(
		*stdout,
		*aws,
		*project,
		*skipImages,
		*skipTimeseries,
		*awsDescription,
		&wg,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = processFiles(extractors.ExtractData, inputFiles, callback)
	if err != nil {
		log.Fatal(err)
	}
	teardown()
}
