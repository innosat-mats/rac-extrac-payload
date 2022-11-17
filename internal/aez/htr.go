package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var htrTemperatures = [...]float64{
	-55, -50, -45, -40, -35, -30, -25, -20, -15, -10,
	-5, 0, 5, 10, 15, 20, 25, 30, 35, 40,
	45, 50, 55, 60, 65, 70, 75, 80, 85, 90,
	95, 100, 105, 110, 115, 120, 125, 130, 135, 140,
	145, 150, 155,
} // ⁰C
var htrResistances = [...]float64{
	9.630e+05, 6.701e+05, 4.717e+05, 3.365e+05, 2.426e+05,
	1.770e+05, 1.304e+05, 9.707e+04, 7.293e+04, 5.533e+04,
	4.232e+04, 3.265e+04, 2.539e+04, 1.990e+04, 1.571e+04,
	1.249e+04, 1.000e+04, 8.057e+03, 6.531e+03, 5.327e+03,
	4.369e+03, 3.603e+03, 2.986e+03, 2.488e+03, 2.083e+03,
	1.752e+03, 1.481e+03, 1.258e+03, 1.072e+03, 9.177e+02,
	7.885e+02, 6.800e+02, 5.886e+02, 5.112e+02, 4.454e+02,
	3.893e+02, 3.417e+02, 3.009e+02, 2.654e+02, 2.348e+02,
	2.083e+02, 1.853e+02, 1.653e+02,
} // Ohm

type htr uint16

func (data *htr) voltage() float64 {
	return voltageConstant * float64(*data)
}

func (data *htr) resistance() float64 {
	return 3.3*3900/data.voltage() - 3900
}

func (data *htr) temperature() (float64, error) {
	return Interpolate(
		data.resistance(),
		htrResistances[:],
		htrTemperatures[:],
	)
}

// HTR housekeeping report returns data on all heater regulators.
type HTR struct {
	HTR1A  htr // Heater 1 Temperature sense A 0..4095
	HTR1B  htr // Heater 1 Temperature sense B 0..4095
	HTR1OD htr // Heater 1 Output Drive setting 0..4095
	HTR2A  htr
	HTR2B  htr
	HTR2OD htr
	HTR7A  htr
	HTR7B  htr
	HTR7OD htr
	HTR8A  htr
	HTR8B  htr
	HTR8OD htr
}

// HTRReport housekeeping report returns data on all heater regulators in useful units.
type HTRReport struct {
	HTR1A    float64 `parquet:"HTR1A"`  // Heater 1 Temperature sense A ⁰C
	HTR1B    float64 `parquet:"HTR1B"`  // Heater 1 Temperature sense B ⁰C
	HTR1OD   float64 `parquet:"HTR1OD"` // Heater 1 Output Drive setting voltage
	HTR2A    float64 `parquet:"HTR2A"`
	HTR2B    float64 `parquet:"HTR2B"`
	HTR2OD   float64 `parquet:"HTR2OD"`
	HTR7A    float64 `parquet:"HTR7A"`
	HTR7B    float64 `parquet:"HTR7B"`
	HTR7OD   float64 `parquet:"HTR7OD"`
	HTR8A    float64 `parquet:"HTR8A"`
	HTR8B    float64 `parquet:"HTR8B"`
	HTR8OD   float64 `parquet:"HTR8OD"`
	WARNINGS []error `parquet:"Warnings"`
}

// NewHTR reads an HTR from buffer
func NewHTR(buf io.Reader) (*HTR, error) {
	htr := HTR{}
	err := binary.Read(buf, binary.LittleEndian, &htr)
	return &htr, err
}

// Report returns a HTRReport with useful units
func (htr *HTR) Report() HTRReport {
	temp1a, err1a := htr.HTR1A.temperature()
	temp1b, err1b := htr.HTR1B.temperature()
	temp2a, err2a := htr.HTR2A.temperature()
	temp2b, err2b := htr.HTR2B.temperature()
	temp7a, err7a := htr.HTR7A.temperature()
	temp7b, err7b := htr.HTR7B.temperature()
	temp8a, err8a := htr.HTR8A.temperature()
	temp8b, err8b := htr.HTR8B.temperature()
	var warnings []error
	if err1a != nil {
		warning := fmt.Errorf("HTR1A: %v", err1a.Error())
		warnings = append(warnings, warning)
	}
	if err1b != nil {
		warning := fmt.Errorf("HTR1B: %v", err1b.Error())
		warnings = append(warnings, warning)
	}
	if err2a != nil {
		warning := fmt.Errorf("HTR2A: %v", err2a.Error())
		warnings = append(warnings, warning)
	}
	if err2b != nil {
		warning := fmt.Errorf("HTR2B: %v", err2b.Error())
		warnings = append(warnings, warning)
	}
	if err7a != nil {
		warning := fmt.Errorf("HTR7A: %v", err7a.Error())
		warnings = append(warnings, warning)
	}
	if err7b != nil {
		warning := fmt.Errorf("HTR7B: %v", err7b.Error())
		warnings = append(warnings, warning)
	}
	if err8a != nil {
		warning := fmt.Errorf("HTR8A: %v", err8a.Error())
		warnings = append(warnings, warning)
	}
	if err8b != nil {
		warning := fmt.Errorf("HTR8B: %v", err8b.Error())
		warnings = append(warnings, warning)
	}
	return HTRReport{
		HTR1A:    temp1a,
		HTR1B:    temp1b,
		HTR1OD:   htr.HTR1OD.voltage(),
		HTR2A:    temp2a,
		HTR2B:    temp2b,
		HTR2OD:   htr.HTR2OD.voltage(),
		HTR7A:    temp7a,
		HTR7B:    temp7b,
		HTR7OD:   htr.HTR7OD.voltage(),
		HTR8A:    temp8a,
		HTR8B:    temp8b,
		HTR8OD:   htr.HTR8OD.voltage(),
		WARNINGS: warnings,
	}
}

// CSVHeaders returns the field names
func (htr *HTR) CSVHeaders() []string {
	return csvHeader(htr.Report())
}

// CSVRow returns the field values
func (htr *HTR) CSVRow() []string {
	val := reflect.Indirect(reflect.ValueOf(htr.Report()))
	values := make([]string, val.NumField())
	t := val.Type()
	for i := range values {
		valueField := val.Field(i)
		if t.Field(i).Name == "WARNINGS" {
			if valueField.Len() == 0 {
				values[i] = ""
			} else {
				var errs = make([]string, valueField.Len())
				for j, l := 0, valueField.Len(); j < l; j++ {
					errs[j] = fmt.Sprintf("%v", valueField.Index(j).Elem())
				}
				values[i] = strings.Join(errs, "|")
			}
		} else {
			values[i] = fmt.Sprintf("%v", valueField.Float())
		}

	}
	return values
}

// CSVSpecifications returns the specs used in creating the struct
func (htr *HTR) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// HTRParquet holds the parquet representation of the HTR
type HTRParquet HTRReport

// GetParquet returns the parquet representation of the HTR
func (htr *HTR) GetParquet() HTRParquet {
	return HTRParquet(htr.Report())
}
