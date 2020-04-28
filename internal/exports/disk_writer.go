package exports

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
)

type csvFile struct {
	File    *os.File
	Writer  *csv.Writer
	HasSpec bool
	HasHead bool
}

type csvOutput interface {
	close()
	setSpecifications(specs []string)
	setHeaderRow(columns []string)
	writeData(data []string)
}

func (csv *csvFile) close() {
	csv.Writer.Flush()
	csv.File.Close()
}

func (csv *csvFile) setSpecifications(specs []string) {
	if csv.HasSpec {
		log.Fatal("Specifications already set for csv output")
	}
	csv.Writer.Write(specs)
	csv.HasSpec = true
}

func (csv *csvFile) setHeaderRow(columns []string) {
	if !csv.HasSpec {
		log.Fatal("Must first supply specifications for csv output")
	}
	if csv.HasHead {
		log.Fatal("Header row already set for csv output")
	}
	csv.Writer.Write(columns)
	csv.HasHead = true
}

func (csv *csvFile) writeData(data []string) {
	if !csv.HasSpec || !csv.HasHead {
		log.Fatal("Specifications and/or Headers missing for csv output")
	}
	csv.Writer.Write(data)
}

func csvOutputFactory(dir string, originName string, packetType string, pkg *ExportablePackage) csvOutput {
	nameParts := strings.Split(filepath.Base(originName), ".")
	name := strings.Join(nameParts[:len(nameParts)-1], ".") + "_" + packetType + ".csv"
	outPath := filepath.Join(
		dir,
		name,
	)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create output directory '%v'", dir)
	}
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Could not create output file '%v'", outPath)
	}
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
) (Callback, CallbackTeardown) {
	var currentOrigin string = ""
	var htrOut csvOutput = nil
	var pwrOut csvOutput = nil
	var cpruOut csvOutput = nil
	var statOut csvOutput = nil

	callback := func(pkg ExportablePackage) {
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
		if pkg.AEZData() == nil {
			return
		}
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
