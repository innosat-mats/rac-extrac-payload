package aez

import (
	"encoding/binary"
	"io"
)

//PWR structure 18 octext
type PWR struct {
	PWRT    uint16 // Temp. sense 0..4095
	PWRP32V uint16 // +32V voltage sense 0..4095
	PWRP32C uint16 // +32V current sense 0..4095
	PWRP16V uint16 // +16V voltage sense 0..4095
	PWRP16C uint16 // +16V current sense 0..4095
	PWRM16V uint16 // -16V voltage sense 0..4095
	PWRM16C uint16 // -16V current sense 0..4095
	PWRP3V3 uint16 // +3V3 voltage sense 0..4095
	PWRP3C3 uint16 // +3V3 current sense 0..4095
}

//PWRReport structure in useful units
type PWRReport struct {
	PWRT    float64 // Temp. sense voltage
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
		PWRT:    pwrVoltageADC(pwr.PWRT),
		PWRP32V: 21 * pwrVoltageADC(pwr.PWRP32V),
		PWRP16V: 11 * pwrVoltageADC(pwr.PWRP16V),
		PWRM16V: -10 * pwrVoltageADC(pwr.PWRM16V),
		PWRP3V3: 4 * pwrVoltageADC(pwr.PWRP3V3),
		PWRP32C: 10.1 / 100 * pwrVoltageADC(pwr.PWRP32C),
		PWRP16C: 10.1 / 5 * pwrVoltageADC(pwr.PWRP16C),
		PWRM16C: 10.1 / 100 * pwrVoltageADC(pwr.PWRM16C),
		PWRP3C3: 10.1 / 20 * pwrVoltageADC(pwr.PWRP3C3),
	}
}
