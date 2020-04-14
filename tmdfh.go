package main

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

func (h *TMDataFieldHeader) read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, h)
}

func (h *TMDataFieldHeader) getPUS() uint8 {
	return (h.PUS << 1) >> 5
}

func (h *TMDataFieldHeader) getServiceType() uint8 {
	return h.ServiceType
}

func (h *TMDataFieldHeader) getServiceSubType() uint8 {
	return h.ServiceSubType
}

func (h *TMDataFieldHeader) getTime() uint32 {
	return h.CUCTimeSeconds
}
