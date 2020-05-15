package exports

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

//CsvFile gives easy access for csv writing
type CsvFile struct {
	file     io.Writer
	writer   *csv.Writer
	Name     string
	HasSpec  bool
	HasHead  bool
	NHeaders int
}

// NewCSVFile returns a CsvFile
func NewCSVFile(out io.Writer, name string) CsvFile {
	return CsvFile{file: out, writer: csv.NewWriter(out), Name: name}
}

// CsvFileWriter implements ease of use writing functions
type CsvFileWriter interface {
	Close()
	SetSpecifications(specs []string) error
	SetHeaderRow(columns []string) error
	WriteData(data []string) error
}

// Close flushes and closes underlying file if any
func (csv *CsvFile) Close() {
	csv.writer.Flush()
	switch csv.file.(type) {
	case *os.File:
		f, ok := csv.file.(*os.File)
		if ok {
			f.Close()
		}
	}
}

// SetSpecifications writes specifications, only allows once
func (csv *CsvFile) SetSpecifications(specs []string) error {
	if csv.HasSpec {
		return fmt.Errorf("Specifications already set for csv output %v", csv.Name)
	}
	csv.writer.Write(specs)
	csv.HasSpec = true
	return nil
}

// SetHeaderRow writes header, only allows once and if specifications previously written
func (csv *CsvFile) SetHeaderRow(columns []string) error {
	if !csv.HasSpec {
		return fmt.Errorf("Must first supply specifications for csv output %v", csv.Name)
	}
	if csv.HasHead {
		return fmt.Errorf("Header row already set for csv output %v", csv.Name)
	}
	csv.writer.Write(columns)
	csv.HasHead = true
	csv.NHeaders = len(columns)
	return nil
}

// WriteData writes a data row, only allows if headers have been written
func (csv *CsvFile) WriteData(data []string) error {
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
	csv.writer.Write(data)
	return nil
}
