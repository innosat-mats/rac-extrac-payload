package exports

import (
	"encoding/csv"
	"fmt"
	"os"
)

type csvFile struct {
	File     *os.File
	Writer   *csv.Writer
	HasSpec  bool
	HasHead  bool
	NHeaders int
}

type csvOutput interface {
	close()
	setSpecifications(specs []string) error
	setHeaderRow(columns []string) error
	writeData(data []string) error
}

func (csv *csvFile) close() {
	csv.Writer.Flush()
	csv.File.Close()
}

func (csv *csvFile) setSpecifications(specs []string) error {
	if csv.HasSpec {
		return fmt.Errorf("Specifications already set for csv output %v", csv.File.Name())
	}
	csv.Writer.Write(specs)
	csv.HasSpec = true
	return nil
}

func (csv *csvFile) setHeaderRow(columns []string) error {
	if !csv.HasSpec {
		return fmt.Errorf("Must first supply specifications for csv output %v", csv.File.Name())
	}
	if csv.HasHead {
		return fmt.Errorf("Header row already set for csv output %v", csv.File.Name())
	}
	csv.Writer.Write(columns)
	csv.HasHead = true
	csv.NHeaders = len(columns)
	return nil
}

func (csv *csvFile) writeData(data []string) error {
	if !csv.HasSpec || !csv.HasHead {
		return fmt.Errorf(
			"Specifications and/or Headers missing for csv output %v",
			csv.File.Name(),
		)
	}
	if csv.NHeaders != len(data) {
		return fmt.Errorf(
			"Irregular column width, expected %v columns but got %v for csv output %v",
			csv.NHeaders,
			len(data),
			csv.File.Name(),
		)
	}
	csv.Writer.Write(data)
	return nil
}
