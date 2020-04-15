package aez

import (
	"encoding/binary"
	"io"
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
