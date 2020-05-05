package exports

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func csvName(dir string, originName string, packetType string) string {
	name := fmt.Sprintf(
		"%v_%v.csv",
		strings.TrimSuffix(filepath.Base(originName), filepath.Ext(originName)),
		packetType,
	)
	return filepath.Join(dir, name)
}

func csvOutputFactory(dir string, originName string, packetType string, pkg *common.ExportablePackage) (csvOutput, error) {
	outPath := csvName(dir, originName, packetType)

	// Create Directory and File
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Could not create output directory '%v'", dir)
	}
	out, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create output file '%v'", outPath)
	}

	// Make a csvFile and produce specs and header row
	csvFile := csvFile{File: out, Writer: csv.NewWriter(out)}
	err = csvFile.setSpecifications((*pkg).CSVSpecifications())
	if err != nil {
		return nil, err
	}
	err = csvFile.setHeaderRow((*pkg).CSVHeaders())
	if err != nil {
		return nil, err
	}
	return &csvFile, nil
}

// DiskCallbackFactory returns a callback for disk writes
func DiskCallbackFactory(
	output string,
	writeImages bool,
	writeTimeseries bool,
) (common.Callback, common.CallbackTeardown) {
	var currentOrigin string = ""
	var err error
	var htrOut csvOutput = nil
	var pwrOut csvOutput = nil
	var cpruOut csvOutput = nil
	var statOut csvOutput = nil

	callback := func(pkg common.ExportablePackage) {
		if writeImages {
			// TODO: Add image writing here
		}

		if !writeTimeseries {
			return
		}
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
				statOut, err = csvOutputFactory(output, currentOrigin, "STAT", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = statOut.writeData(pkg.CSVRow())
		case aez.HTR:
			if htrOut == nil {
				htrOut, err = csvOutputFactory(output, currentOrigin, "HTR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = htrOut.writeData(pkg.CSVRow())
		case aez.PWR:
			if pwrOut == nil {
				pwrOut, err = csvOutputFactory(output, currentOrigin, "PWR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = pwrOut.writeData(pkg.CSVRow())
		case aez.CPRU:
			if cpruOut == nil {
				cpruOut, err = csvOutputFactory(output, currentOrigin, "CPRU", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = cpruOut.writeData(pkg.CSVRow())
		}
		// This error comes from writing a line and most probably would be a column missmatch
		// that means we should be able to continue and just report the error
		if err != nil {
			log.Println(err)
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