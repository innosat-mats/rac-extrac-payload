package aez

import (
	"encoding/binary"
	"io"
)

//HTR housekeeping report returns data on all heater regulators.
type HTR struct {
	HTR1A  uint16 // Heater 1 Temperature sense A 0..4095
	HTR1B  uint16 // Heater 1 Temperature sense B 0..4095
	HTR1OD uint16 // Heater 1 Output Drive setting 0..4095
	HTR2A  uint16
	HTR2B  uint16
	HTR2OD uint16
	HTR3A  uint16
	HTR3B  uint16
	HTR3OD uint16
	HTR4A  uint16
	HTR4B  uint16
	HTR4OD uint16
	HTR5A  uint16
	HTR5B  uint16
	HTR5OD uint16
	HTR6A  uint16
	HTR6B  uint16
	HTR6OD uint16
	HTR7A  uint16
	HTR7B  uint16
	HTR7OD uint16
	HTR8A  uint16
	HTR8B  uint16
	HTR8OD uint16
}

//HTRReport housekeeping report returns data on all heater regulators in useful units.
type HTRReport struct {
	HTR1A  float64 // Heater 1 Temperature sense A voltage
	HTR1B  float64 // Heater 1 Temperature sense B voltage
	HTR1OD float64 // Heater 1 Output Drive setting voltage
	HTR2A  float64
	HTR2B  float64
	HTR2OD float64
	HTR3A  float64
	HTR3B  float64
	HTR3OD float64
	HTR4A  float64
	HTR4B  float64
	HTR4OD float64
	HTR5A  float64
	HTR5B  float64
	HTR5OD float64
	HTR6A  float64
	HTR6B  float64
	HTR6OD float64
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

func htrVoltage(data uint16) float64 {
	return voltageConstant * float64(data)
}

// Report returns a HTRReport with useful units
func (htr *HTR) Report() HTRReport {
	return HTRReport{
		HTR1A:  htrVoltage(htr.HTR1A),
		HTR1B:  htrVoltage(htr.HTR1B),
		HTR1OD: htrVoltage(htr.HTR1OD),
		HTR2A:  htrVoltage(htr.HTR2A),
		HTR2B:  htrVoltage(htr.HTR2B),
		HTR2OD: htrVoltage(htr.HTR2OD),
		HTR3A:  htrVoltage(htr.HTR3A),
		HTR3B:  htrVoltage(htr.HTR3B),
		HTR3OD: htrVoltage(htr.HTR3OD),
		HTR4A:  htrVoltage(htr.HTR4A),
		HTR4B:  htrVoltage(htr.HTR4B),
		HTR4OD: htrVoltage(htr.HTR4OD),
		HTR5A:  htrVoltage(htr.HTR5A),
		HTR5B:  htrVoltage(htr.HTR5B),
		HTR5OD: htrVoltage(htr.HTR5OD),
		HTR6A:  htrVoltage(htr.HTR6A),
		HTR6B:  htrVoltage(htr.HTR6B),
		HTR6OD: htrVoltage(htr.HTR6OD),
		HTR7A:  htrVoltage(htr.HTR7A),
		HTR7B:  htrVoltage(htr.HTR7B),
		HTR7OD: htrVoltage(htr.HTR7OD),
		HTR8A:  htrVoltage(htr.HTR8A),
		HTR8B:  htrVoltage(htr.HTR8B),
		HTR8OD: htrVoltage(htr.HTR8OD),
	}
}
