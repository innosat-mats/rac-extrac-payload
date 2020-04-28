package exports

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
)

func flattenDeepInternal(args []string, v reflect.Value) []string {

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			args = flattenDeepInternal(args, v.Index(i))
		}
	} else {

		args = append(args, v.String())
	}

	return args
}

func flatten(args [][]string) []string {
	return flattenDeepInternal(nil, reflect.ValueOf(args))
}

type csvFile struct {
	File    *os.File
	Writer  *csv.Writer
	HasSpec bool
	HasHead bool
}

type csvOutput interface {
	close()
	setSpecifications(specs []string)
	setHeaderRow(columns ...[]string)
	writeData(data ...[]string)
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

func (csv *csvFile) setHeaderRow(columns ...[]string) {
	if !csv.HasSpec {
		log.Fatal("Must first supply specifications for csv output")
	}
	if csv.HasHead {
		log.Fatal("Header row already set for csv output")
	}
	csv.Writer.Write(flatten(columns))
	csv.HasHead = true
}

func (csv *csvFile) writeData(data ...[]string) {
	csv.Writer.Write(flatten(data))
}

func csvOutputFactory(dir string, originName string, packetType string) csvOutput {
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
	return &csvFile{File: out, Writer: csv.NewWriter(out)}
}

// DiskCallbackFactory returns a callback for disk writes
func DiskCallbackFactory(
	output string,
	writeImages bool,
	writeTimeseries bool,
) Callback {
	var currentOrigin string = ""
	var htrOut csvOutput = nil
	var pwrOut csvOutput = nil
	return func(pkg ExportablePackage) {
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
			currentOrigin = pkg.OriginName()
		}
		if pkg.AEZData() == nil {
			return
		}
		switch pkg.AEZData().(type) {
		case aez.HTR:
			if htrOut == nil {
				htrOut = csvOutputFactory(output, currentOrigin, "HTR")
				defer htrOut.close()
				htrOut.setSpecifications(pkg.CSVSpecifications())
				htrOut.setHeaderRow(pkg.CSVHeaders())
			}
			htrOut.writeData(pkg.CSVRow())
		}
	}
}
