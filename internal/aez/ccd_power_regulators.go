package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

//CPRU structure
type CPRU struct {
	VGATE0 uint16 // CCD0 Gate Voltage 0..4095
	VSUBS0 uint16 // CCD0 Substrate Voltage 0..4095
	VRD0   uint16 // CCD0 Reset transistor Drain Voltage 0..4095
	VOD0   uint16 // CCD0 Output Drain Voltage 0..4095
	VGATE1 uint16 // CCD1 Gate Voltage 0..4095
	VSUBS1 uint16 // CCD1 Substrate Voltage 0..4095
	VRD1   uint16 // CCD1 Reset transistor Drain Voltage 0..4095
	VOD1   uint16 // CCD1 Output Drain Voltage 0..4095
	VGATE2 uint16 // CCD2 Gate Voltage 0..4095
	VSUBS2 uint16 // CCD2 Substrate Voltage 0..4095
	VRD2   uint16 // CCD2 Reset transistor Drain Voltage 0..4095
	VOD2   uint16 // CCD2 Output Drain Voltage 0..4095
	VGATE3 uint16 // CCD3 Gate Voltage 0..4095
	VSUBS3 uint16 // CCD3 Substrate Voltage 0..4095
	VRD3   uint16 // CCD3 Reset transistor Drain Voltage 0..4095
	VOD3   uint16 // CCD3 Output Drain Voltage 0..4095
}

// Read CRPU
func (cpru *CPRU) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, cpru)
}

var voltageConstant float64 = 2.5 / (math.Pow(2, 12) - 1)

func gateVoltage(data uint16) float64 {
	return voltageConstant * float64(data) * 10
}

func subsVoltage(data uint16) float64 {
	return voltageConstant * float64(data) * 11 / 1.5
}

func rdVoltage(data uint16) float64 {
	return voltageConstant * float64(data) * 17 / 1.5
}

func odVoltage(data uint16) float64 {
	return voltageConstant * float64(data) * 32 / 1.5
}

// CPRUField is enum like positions of CPRU fields in the struct
type CPRUField int

// Enum-like positions of fileds in the CPRU struct
const (
	VGATE0 CPRUField = iota
	VSUBS0
	VRD0
	VOD0
	VGATE1
	VSUBS1
	VRD1
	VOD1
	VGATE2
	VSUBS2
	VRD2
	VOD2
	VGATE3
	VSUBS3
	VRD3
	VOD3
)

func (cpru *CPRU) extractField(field CPRUField) uint16 {
	reflectedValue := reflect.ValueOf(cpru)
	reflectedField := reflect.Indirect(reflectedValue).Field(int(field))
	return uint16(reflectedField.Uint())
}

// Voltage returns the voltage for the specified CPRUField
func (cpru *CPRU) Voltage(field CPRUField) (float64, error) {
	switch field {
	case VGATE0, VGATE1, VGATE2, VGATE3:
		return gateVoltage(cpru.extractField(field)), nil
	case VSUBS0, VSUBS1, VSUBS2, VSUBS3:
		return subsVoltage(cpru.extractField(field)), nil
	case VRD0, VRD1, VRD2, VRD3:
		return rdVoltage(cpru.extractField(field)), nil
	case VOD0, VOD1, VOD2, VOD3:
		return odVoltage(cpru.extractField(field)), nil
	}
	return 0, fmt.Errorf("Unknown CPRUField %d", field)
}
