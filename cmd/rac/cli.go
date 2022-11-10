package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/exports"
	"github.com/innosat-mats/rac-extract-payload/internal/extractors"
)

// Version is the version of the source code
var Version string

// Head is the short commit id of head
var Head string

// Buildtime is the time of the build
var Buildtime string

var skipImages *bool
var skipTimeseries *bool
var project *string
var stdout *bool
var dregsDir *string
var version *bool

// myUsage replaces default usage since it doesn't include information on non-flags
func myUsage() {
	fmt.Println("Extracts information from Innosat-MATS rac-files")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTIONS] rac-file ...\n", os.Args[0])
	if len(os.Args) > 2 {
		switch helpSection := strings.ToUpper(os.Args[2]); helpSection {
		case "OUTPUT":
			infoGeneral()
		case "CCD":
			infoCCD()
		case "CPRU":
			infoCPRU()
		case "HTR":
			infoHTR()
		case "PWR":
			infoPWR()
		case "STAT":
			infoSTAT()
		case "TCV":
			infoTCV()
		case "PM":
			infoPM()
		case "MATS", "SPACE", "M.A.T.S.", "SATELLITE":
			infoSpace()
		default:
			fmt.Printf("\nUnrecognized help section %s\n", helpSection)
		}
		return
	}
	flag.PrintDefaults()
	fmt.Printf(
		"\nFor extra information about the output CSV:s type \"%s -help output\"\n",
		os.Args[0],
	)

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
	project string,
	skipImages bool,
	skipTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown, error) {
	if project == "" && !toStdout {
		flag.Usage()
		fmt.Println("\nExpected a project")
		return nil, nil, errors.New("invalid arguments")
	}
	if skipTimeseries && (skipImages || toStdout) {
		fmt.Println("Nothing will be extracted, only validating integrity of rac-file(s)")
	}

	if toStdout {
		callback, teardown := exports.StdoutCallbackFactory(os.Stdout, !skipTimeseries)
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
	dregs extractors.Dregs,
	callback common.Callback,
) error {
	batch := make([]extractors.StreamBatch, len(inputFiles))
	for n, filename := range inputFiles {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		batch[n] = extractors.StreamBatch{
			Buf: f,
			Origin: &common.OriginDescription{
				Name:           filename,
				ProcessingDate: time.Now(),
			},
		}

	}
	extractor(callback, dregs, batch...)
	return nil
}

func init() {
	common.Version = Version
	common.Head = Head
	common.Buildtime = Buildtime

	skipImages = flag.Bool("skip-images", false, "Extract images from rac-files.\n(Default: false)")
	skipTimeseries = flag.Bool(
		"skip-timeseries",
		false,
		"Extract timeseries from rac-files.\n(Default: false)",
	)
	project = flag.String(
		"project",
		"",
		"Name for experiments, when outputting to disk a directory will be created with this name.",
	)
	stdout = flag.Bool(
		"stdout",
		false,
		"Output to standard out instead of to disk (only timeseries)\n(Default: false)",
	)
	dregsDir = flag.String(
		"dregs",
		"",
		"Path to directory where to find and write dregs files for multi packet continuation. Directory will be created if non-existent. If empty dregs will be skipped.",
	)
	version = flag.Bool(
		"version",
		false,
		"Only display current version of the program",
	)

	flag.Usage = myUsage
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	if *version {
		fmt.Println("Version", Version, "Commit", Head, "@", Buildtime)
		return
	}

	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		flag.Usage()
		log.Fatal("No rac-files supplied")
	}
	callback, teardown, err := getCallback(
		*stdout,
		*project,
		*skipImages,
		*skipTimeseries,
		&wg,
	)
	if err != nil {
		log.Fatal(err)
	}
	dregs := extractors.Dregs{
		Path:    *dregsDir,
		MaxDiff: extractors.MaxDeviationNanos,
	}
	err = processFiles(extractors.ExtractData, inputFiles, dregs, callback)
	if err != nil {
		log.Fatal(err)
	}
	teardown()
}
