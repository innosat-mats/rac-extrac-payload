package aez

import (
	"encoding/binary"
	"io"
	"math"
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

//CPRUReport structure
type CPRUReport struct {
	VGATE0 float64 // CCD0 Gate Voltage
	VSUBS0 float64 // CCD0 Substrate Voltage
	VRD0   float64 // CCD0 Reset transistor Drain Voltage
	VOD0   float64 // CCD0 Output Drain Voltage
	VGATE1 float64 // CCD1 Gate Voltage
	VSUBS1 float64 // CCD1 Substrate Voltage
	VRD1   float64 // CCD1 Reset transistor Drain Voltage
	VOD1   float64 // CCD1 Output Drain Voltage
	VGATE2 float64 // CCD2 Gate Voltage
	VSUBS2 float64 // CCD2 Substrate Voltage
	VRD2   float64 // CCD2 Reset transistor Drain Voltage
	VOD2   float64 // CCD2 Output Drain Voltage
	VGATE3 float64 // CCD3 Gate Voltage
	VSUBS3 float64 // CCD3 Substrate Voltage
	VRD3   float64 // CCD3 Reset transistor Drain Voltage
	VOD3   float64 // CCD3 Output Drain Voltage
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

// Report transforms CPRU data to useful units
func (cpru *CPRU) Report() CPRUReport {
	return CPRUReport{
		VGATE0: gateVoltage(cpru.VGATE0),
		VSUBS0: subsVoltage(cpru.VSUBS0),
		VRD0:   rdVoltage(cpru.VRD0),
		VOD0:   odVoltage(cpru.VOD0),
		VGATE1: gateVoltage(cpru.VGATE1),
		VSUBS1: subsVoltage(cpru.VSUBS1),
		VRD1:   rdVoltage(cpru.VRD1),
		VOD1:   odVoltage(cpru.VOD1),
		VGATE2: gateVoltage(cpru.VGATE2),
		VSUBS2: subsVoltage(cpru.VSUBS2),
		VRD2:   rdVoltage(cpru.VRD2),
		VOD2:   odVoltage(cpru.VOD2),
		VGATE3: gateVoltage(cpru.VGATE3),
		VSUBS3: subsVoltage(cpru.VSUBS3),
		VRD3:   rdVoltage(cpru.VRD3),
		VOD3:   odVoltage(cpru.VOD3),
	}
}
