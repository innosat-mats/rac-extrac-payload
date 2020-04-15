package innosat

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
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
func (h *TMDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, h)
}

// PUSVersion ...
func (h *TMDataFieldHeader) PUSVersion() uint8 {
	return (h.PUS << 1) >> 5
}

// Time returns the telemetry data time in UTC
func (h *TMDataFieldHeader) Time() time.Time {
	return ccsds.UnsegmentedTimeDate(h.CUCTimeSeconds, h.CUCTimeFraction)
}

// Nanoseconds returns the telemetry data time in nanoseconds since its epoch
func (h *TMDataFieldHeader) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(h.CUCTimeSeconds, h.CUCTimeFraction)
}
