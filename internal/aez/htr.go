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

// Read ..
func (s *HTR) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, s)
}
