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
func (tmdfh *TMDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, tmdfh)
}

// PUSVersion ...
func (tmdfh TMDataFieldHeader) PUSVersion() uint8 {
	return (tmdfh.PUS << 1) >> 5
}

// Time returns the telemetry data time in UTC
func (tmdfh *TMDataFieldHeader) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(tmdfh.CUCTimeSeconds, tmdfh.CUCTimeFraction, epoch)
}

// Nanoseconds returns the telemetry data time in nanoseconds since its epoch
func (tmdfh *TMDataFieldHeader) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(tmdfh.CUCTimeSeconds, tmdfh.CUCTimeFraction)
}

// IsHousekeeping returns if payload contains housekeeping data
func (tmdfh *TMDataFieldHeader) IsHousekeeping() bool {
	return tmdfh.ServiceType == 3 && tmdfh.ServiceSubType == 25
}
