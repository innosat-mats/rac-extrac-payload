package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
)

type gate uint16

var voltageConstant float64 = 2.5 / (math.Pow(2, 12) - 1)

func (data gate) voltage() float64 {
	return voltageConstant * float64(data) * 10
}

type subs uint16

func (data subs) voltage() float64 {
	return voltageConstant * float64(data) * 11 / 1.5
}

type rd uint16

func (data rd) voltage() float64 {
	return voltageConstant * float64(data) * 17 / 1.5
}

type od uint16

func (data od) voltage() float64 {
	return voltageConstant * float64(data) * 32 / 1.5
}

type cpruStat uint8

func (stat cpruStat) overvoltageFault(ccd uint8) bool {
	if ccd > 3 {
		log.Fatalf("overvoltageFault got illegal ccd %v", ccd)
	}
	var mask uint8 = 0x80 >> ccd
	return uint8(stat)&mask != 0
}

func (stat cpruStat) powerEnabled(ccd uint8) bool {
	if ccd > 3 {
		log.Fatalf("powerEnabled got illegal ccd %v", ccd)
	}
	var mask uint8 = 0x08 >> ccd
	return uint8(stat)&mask != 0
}

//CPRU structure
type CPRU struct {
	STAT cpruStat // CPRU/CRB power status
	// CCD overvoltage fault, one bit per CCD. Bit [7..4]
	// CCD power enabled, one bit per CCD. Bit [3..0]
	VGATE0 gate // CCD0 Gate Voltage 0..4095
	VSUBS0 subs // CCD0 Substrate Voltage 0..4095
	VRD0   rd   // CCD0 Reset transistor Drain Voltage 0..4095
	VOD0   od   // CCD0 Output Drain Voltage 0..4095
	VGATE1 gate // CCD1 Gate Voltage 0..4095
	VSUBS1 subs // CCD1 Substrate Voltage 0..4095
	VRD1   rd   // CCD1 Reset transistor Drain Voltage 0..4095
	VOD1   od   // CCD1 Output Drain Voltage 0..4095
	VGATE2 gate // CCD2 Gate Voltage 0..4095
	VSUBS2 subs // CCD2 Substrate Voltage 0..4095
	VRD2   rd   // CCD2 Reset transistor Drain Voltage 0..4095
	VOD2   od   // CCD2 Output Drain Voltage 0..4095
	VGATE3 gate // CCD3 Gate Voltage 0..4095
	VSUBS3 subs // CCD3 Substrate Voltage 0..4095
	VRD3   rd   // CCD3 Reset transistor Drain Voltage 0..4095
	VOD3   od   // CCD3 Output Drain Voltage 0..4095
}

//CPRUReport structure
type CPRUReport struct {
	VGATE0       float64 // CCD0 Gate Voltage
	VSUBS0       float64 // CCD0 Substrate Voltage
	VRD0         float64 // CCD0 Reset transistor Drain Voltage
	VOD0         float64 // CCD0 Output Drain Voltage
	Overvoltage0 bool    // CCD0 overvoltage fault
	Power0       bool    // CCD0 overvoltage fault
	VGATE1       float64 // CCD1 Gate Voltage
	VSUBS1       float64 // CCD1 Substrate Voltage
	VRD1         float64 // CCD1 Reset transistor Drain Voltage
	VOD1         float64 // CCD1 Output Drain Voltage
	Overvoltage1 bool    // CCD1 overvoltage fault
	Power1       bool    // CCD1 overvoltage fault
	VGATE2       float64 // CCD2 Gate Voltage
	VSUBS2       float64 // CCD2 Substrate Voltage
	VRD2         float64 // CCD2 Reset transistor Drain Voltage
	VOD2         float64 // CCD2 Output Drain Voltage
	Overvoltage2 bool    // CCD2 overvoltage fault
	Power2       bool    // CCD2 overvoltage fault
	VGATE3       float64 // CCD3 Gate Voltage
	VSUBS3       float64 // CCD3 Substrate Voltage
	VRD3         float64 // CCD3 Reset transistor Drain Voltage
	VOD3         float64 // CCD3 Output Drain Voltage
	Overvoltage3 bool    // CCD3 overvoltage fault
	Power3       bool    // CCD3 overvoltage fault

}

// Read CRPU
func (cpru *CPRU) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, cpru)
}

// Report transforms CPRU data to useful units
func (cpru *CPRU) Report() CPRUReport {
	return CPRUReport{
		VGATE0:       cpru.VGATE0.voltage(),
		VSUBS0:       cpru.VSUBS0.voltage(),
		VRD0:         cpru.VRD0.voltage(),
		VOD0:         cpru.VOD0.voltage(),
		Overvoltage0: cpru.STAT.overvoltageFault(0),
		Power0:       cpru.STAT.powerEnabled(0),
		VGATE1:       cpru.VGATE1.voltage(),
		VSUBS1:       cpru.VSUBS1.voltage(),
		VRD1:         cpru.VRD1.voltage(),
		VOD1:         cpru.VOD1.voltage(),
		Overvoltage1: cpru.STAT.overvoltageFault(1),
		Power1:       cpru.STAT.powerEnabled(1),
		VGATE2:       cpru.VGATE2.voltage(),
		VSUBS2:       cpru.VSUBS2.voltage(),
		VRD2:         cpru.VRD2.voltage(),
		VOD2:         cpru.VOD2.voltage(),
		Overvoltage2: cpru.STAT.overvoltageFault(2),
		Power2:       cpru.STAT.powerEnabled(2),
		VGATE3:       cpru.VGATE3.voltage(),
		VSUBS3:       cpru.VSUBS3.voltage(),
		VRD3:         cpru.VRD3.voltage(),
		VOD3:         cpru.VOD3.voltage(),
		Overvoltage3: cpru.STAT.overvoltageFault(3),
		Power3:       cpru.STAT.powerEnabled(3),
	}
}

//CSVSpecifications returns the specs used in creating the struct
func (cpru CPRU) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

//CSVHeaders returns the field names
func (cpru CPRU) CSVHeaders() []string {
	return csvHeader(cpru.Report())
}

//CSVRow returns the field values
func (cpru CPRU) CSVRow() []string {
	val := reflect.Indirect(reflect.ValueOf(cpru.Report()))
	values := make([]string, val.NumField())
	for i := range values {
		valueField := val.Field(i)
		switch valueField.Interface().(type) {
		case bool:
			values[i] = fmt.Sprintf("%v", valueField.Bool())
		default:
			values[i] = fmt.Sprintf("%v", valueField.Float())
		}
	}
	return values
}
