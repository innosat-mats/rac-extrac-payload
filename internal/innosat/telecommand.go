package innosat

import (
	"encoding/binary"
	"io"
)

//TCDataFieldHeader (9 octets)
type TCDataFieldHeader struct {
	PUS            uint8
	ServiceType    uint8
	ServiceSubType uint8
}

// Read TCDataFieldHeader
func (h *TCDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, &h)
}

//PUSVersion ...
func (h *TCDataFieldHeader) PUSVersion() uint8 {
	return (h.PUS << 1) >> 5
}
