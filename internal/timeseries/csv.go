package timeseries

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

//CSV gives easy access for csv writing
type CSV struct {
	writer    io.Writer
	csvWriter *csv.Writer
	Name      string
	HasSpec   bool
	HasHead   bool
	NHeaders  int
}

// NewCSV returns a Timeseries CSV
func NewCSV(out io.Writer, name string) CSVWriter {
	return &CSV{writer: out, csvWriter: csv.NewWriter(out), Name: name}
}

// CSVWriter implements ease of use writing functions
type CSVWriter interface {
	Close()
	SetSpecifications(specs []string) error
	SetHeaderRow(columns []string) error
	WriteData(data []string) error
}

// Close flushes and closes underlying file if any
func (csv *CSV) Close() {
	csv.csvWriter.Flush()
	switch csv.writer.(type) {
	case *os.File:
		f, ok := csv.writer.(*os.File)
		if ok {
			f.Close()
		}
	}
}

// SetSpecifications writes specifications, only allows once
func (csv *CSV) SetSpecifications(specs []string) error {
	if csv.HasSpec {
		return fmt.Errorf("Specifications already set for csv output %v", csv.Name)
	}
	csv.csvWriter.Write(specs)
	csv.HasSpec = true
	return nil
}

// SetHeaderRow writes header, only allows once and if specifications previously written
func (csv *CSV) SetHeaderRow(columns []string) error {
	if !csv.HasSpec {
		return fmt.Errorf("Must first supply specifications for csv output %v", csv.Name)
	}
	if csv.HasHead {
		return fmt.Errorf("Header row already set for csv output %v", csv.Name)
	}
	csv.csvWriter.Write(columns)
	csv.HasHead = true
	csv.NHeaders = len(columns)
	return nil
}

// WriteData writes a data row, only allows if headers have been written
func (csv *CSV) WriteData(data []string) error {
	if !csv.HasSpec || !csv.HasHead {
		return fmt.Errorf(
			"Specifications and/or Headers missing for csv output %v",
			csv.Name,
		)
	}
	if csv.NHeaders != len(data) {
		return fmt.Errorf(
			"Irregular column width, expected %v columns but got %v for csv output %v",
			csv.NHeaders,
			len(data),
			csv.Name,
		)
	}
	csv.csvWriter.Write(data)
	return nil
}
