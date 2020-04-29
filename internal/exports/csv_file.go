package exports

import (
	"encoding/csv"
	"log"
	"os"
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