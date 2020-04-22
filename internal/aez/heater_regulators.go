package aez

import (
	"encoding/binary"
	"io"
)

type htr uint16

func (data htr) voltage() float64 {
	return voltageConstant * float64(data)
}

func (data htr) resistance() float64 {
	return 3.3*3900/data.voltage() - 3900
}

//HTR housekeeping report returns data on all heater regulators.
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

//HTRReport housekeeping report returns data on all heater regulators in useful units.
type HTRReport struct {
	HTR1A  float64 // Heater 1 Temperature sense A voltage
	HTR1B  float64 // Heater 1 Temperature sense B resistance
	HTR1OD float64 // Heater 1 Output Drive setting voltage
	HTR2A  float64
	HTR2B  float64
	HTR2OD float64
	HTR7A  float64
	HTR7B  float64
	HTR7OD float64
	HTR8A  float64
	HTR8B  float64
	HTR8OD float64
}

// Read HTR
func (htr *HTR) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, htr)
}

// Report returns a HTRReport with useful units
func (htr *HTR) Report() HTRReport {
	return HTRReport{
		HTR1A:  htr.HTR1A.voltage(),
		HTR1B:  htr.HTR1B.resistance(),
		HTR1OD: htr.HTR1OD.voltage(),
		HTR2A:  htr.HTR2A.voltage(),
		HTR2B:  htr.HTR2B.resistance(),
		HTR2OD: htr.HTR2OD.voltage(),
		HTR7A:  htr.HTR7A.voltage(),
		HTR7B:  htr.HTR7B.resistance(),
		HTR7OD: htr.HTR7OD.voltage(),
		HTR8A:  htr.HTR8A.voltage(),
		HTR8B:  htr.HTR8B.resistance(),
		HTR8OD: htr.HTR8OD.voltage(),
	}
}
