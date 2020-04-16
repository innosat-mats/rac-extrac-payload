package innosat

import (
	"encoding/binary"
	"io"
)

//TMDataFieldHeader (9 octets)
type TMDataFieldHeader struct {
	PUS             uint8
	ServiceType     uint8
	ServiceSubType  uint8
	CUCTimeSeconds  uint32
	CUCTimeFraction uint16
}

// Read TMDataFieldHeader
func (tmdfh *TMDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, tmdfh)
}

// PUSVersion ...
func (tmdfh *TMDataFieldHeader) PUSVersion() uint8 {
	return (tmdfh.PUS << 1) >> 5
}

// Time ...
func (tmdfh *TMDataFieldHeader) Time() uint32 {
	return tmdfh.CUCTimeSeconds
}
