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
func (tcdfh *TCDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, &tcdfh)
}

//PUSVersion ...
func (tcdfh *TCDataFieldHeader) PUSVersion() uint8 {
	return (tcdfh.PUS << 1) >> 5
}
