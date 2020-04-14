package pus

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

// Read ...
func (h *TMDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, h)
}

// PUSVersion ...
func (h *TMDataFieldHeader) PUSVersion() uint8 {
	return (h.PUS << 1) >> 5
}

// Time ...
func (h *TMDataFieldHeader) Time() uint32 {
	return h.CUCTimeSeconds
}
