package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type pwr uint16

func (data pwr) voltageADC() float64 {
	return voltageConstant * float64(data)
}

type pwrt pwr

func (data pwrt) voltage() float64 {
	return pwr(data).voltageADC()
}

func (data pwrt) resistance() float64 {
	return 3.3*1000/data.voltage() - 1000
}

type pwrp32v pwr

func (data pwrp32v) voltage() float64 {
	return 21 * pwr(data).voltageADC()
}

type pwrp32c pwr

func (data pwrp32c) current() float64 {
	return 10.1 / 100 * pwr(data).voltageADC()
}

type pwrp16v pwr

func (data pwrp16v) voltage() float64 {
	return 11 * pwr(data).voltageADC()
}

type pwrp16c pwr

func (data pwrp16c) current() float64 {
	return 10.1 / 5 * pwr(data).voltageADC()
}

type pwrm16v pwr

func (data pwrm16v) voltage() float64 {
	return -10 * pwr(data).voltageADC()
}

type pwrm16c pwr

func (data pwrm16c) current() float64 {
	return 10.1 / 100 * pwr(data).voltageADC()
}

type pwrp3v3 pwr

func (data pwrp3v3) voltage() float64 {
	return 4 * pwr(data).voltageADC()
}

type pwrp3c3 pwr

func (data pwrp3c3) current() float64 {
	return 10.1 / 20 * pwr(data).voltageADC()
}

//PWR structure 18 octext
type PWR struct {
	PWRT    pwrt    // Temp. sense 0..4095
	PWRP32V pwrp32v // +32V voltage sense 0..4095
	PWRP32C pwrp32c // +32V current sense 0..4095
	PWRP16V pwrp16v // +16V voltage sense 0..4095
	PWRP16C pwrp16c // +16V current sense 0..4095
	PWRM16V pwrm16v // -16V voltage sense 0..4095
	PWRM16C pwrm16c // -16V current sense 0..4095
	PWRP3V3 pwrp3v3 // +3V3 voltage sense 0..4095
	PWRP3C3 pwrp3c3 // +3V3 current sense 0..4095
}

//PWRReport structure in useful units
type PWRReport struct {
	PWRT    float64 // Temp. sense resistance
	PWRP32V float64 // +32V voltage sense voltage
	PWRP32C float64 // +32V current sense current
	PWRP16V float64 // +16V voltage sense voltage
	PWRP16C float64 // +16V current sense current
	PWRM16V float64 // -16V voltage sense voltage
	PWRM16C float64 // -16V current sense current
	PWRP3V3 float64 // +3V3 voltage sense voltage
	PWRP3C3 float64 // +3V3 current sense current
}

// Read PWR
func (pwr *PWR) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, pwr)
}

func pwrVoltageADC(data uint16) float64 {
	return voltageConstant * float64(data)
}

// Report returns a PWRReport with useful units
func (pwr *PWR) Report() PWRReport {
	return PWRReport{
		PWRT:    pwr.PWRT.resistance(),
		PWRP32V: pwr.PWRP32V.voltage(),
		PWRP32C: pwr.PWRP32C.current(),
		PWRP16V: pwr.PWRP16V.voltage(),
		PWRP16C: pwr.PWRP16C.current(),
		PWRM16V: pwr.PWRM16V.voltage(),
		PWRM16C: pwr.PWRM16C.current(),
		PWRP3V3: pwr.PWRP3V3.voltage(),
		PWRP3C3: pwr.PWRP3C3.current(),
	}
}

//CSVSpecifications returns the specs used in creating the struct
func (pwr PWR) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

//CSVHeaders returns the field names
func (pwr PWR) CSVHeaders() []string {
	val := reflect.Indirect(reflect.ValueOf(pwr.Report()))
	t := val.Type()
	fields := make([]string, t.NumField())
	for i := range fields {
		fields[i] = t.Field(i).Name
	}
	return fields
}

//CSVRow returns the field values
func (pwr PWR) CSVRow() []string {
	val := reflect.Indirect(reflect.ValueOf(pwr.Report()))
	values := make([]string, val.NumField())
	for i := range values {
		valueField := val.Field(i)
		values[i] = fmt.Sprintf("%f", valueField.Float())
	}
	return values
}
