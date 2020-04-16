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

// Read PWR
func (pwr *PWR) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, pwr)
}
