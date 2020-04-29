package exports

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func csvName(dir string, originName string, packetType string) string {
	nameParts := strings.Split(filepath.Base(originName), ".")
	name := strings.Join(nameParts[:len(nameParts)-1], ".") + "_" + packetType + ".csv"
	return filepath.Join(dir, name)
}

func csvOutputFactory(dir string, originName string, packetType string, pkg *common.ExportablePackage) csvOutput {
	outPath := csvName(dir, originName, packetType)

	// Create Directory and File
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create output directory '%v'", dir)
	}
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Could not create output file '%v'", outPath)
	}

	// Make a csvFile and produce specs and header row
	file := csvFile{File: out, Writer: csv.NewWriter(out)}
	file.setSpecifications((*pkg).CSVSpecifications())
	file.setHeaderRow((*pkg).CSVHeaders())
	return &file
}

// DiskCallbackFactory returns a callback for disk writes
func DiskCallbackFactory(
	output string,
	writeImages bool,
	writeTimeseries bool,
) (common.Callback, common.CallbackTeardown) {
	var currentOrigin string = ""
	var htrOut csvOutput = nil
	var pwrOut csvOutput = nil
	var cpruOut csvOutput = nil
	var statOut csvOutput = nil

	callback := func(pkg common.ExportablePackage) {
		// Close streams from previous file
		if pkg.OriginName() != currentOrigin {
			if pwrOut != nil {
				pwrOut.close()
				pwrOut = nil
			}
			if htrOut != nil {
				htrOut.close()
				htrOut = nil
			}
			if statOut != nil {
				statOut.close()
				statOut = nil
			}
			if cpruOut != nil {
				cpruOut.close()
				cpruOut = nil
			}
			currentOrigin = pkg.OriginName()
		}

		// We have nowhere to write partial extraction of record so we discard
		if pkg.AEZData() == nil {
			return
		}

		// Write to the dedicated target stream
		switch pkg.AEZData().(type) {
		case aez.STAT:
			if statOut == nil {
				statOut = csvOutputFactory(output, currentOrigin, "STAT", &pkg)
			}
			statOut.writeData(pkg.CSVRow())
		case aez.HTR:
			if htrOut == nil {
				htrOut = csvOutputFactory(output, currentOrigin, "HTR", &pkg)
			}
			htrOut.writeData(pkg.CSVRow())
		case aez.PWR:
			if pwrOut == nil {
				pwrOut = csvOutputFactory(output, currentOrigin, "PWR", &pkg)
			}
			pwrOut.writeData(pkg.CSVRow())
		case aez.CPRU:
			if cpruOut == nil {
				cpruOut = csvOutputFactory(output, currentOrigin, "CPRU", &pkg)
			}
			cpruOut.writeData(pkg.CSVRow())
		}
	}

	teardown := func() {
		if statOut != nil {
			statOut.close()
		}
		if htrOut != nil {
			htrOut.close()
		}
		if pwrOut != nil {
			pwrOut.close()
		}
		if cpruOut != nil {
			cpruOut.close()
		}
	}

	return callback, teardown
}
