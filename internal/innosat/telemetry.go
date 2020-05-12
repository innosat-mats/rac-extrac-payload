package innosat

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

type pus uint8

// Version ...
func (pus pus) Version() uint8 {
	return uint8((pus << 1) >> 5)
}

//TMDataFieldHeader (9 octets)
type TMDataFieldHeader struct {
	PUS             pus
	ServiceType     SourcePackageServiceType
	ServiceSubType  uint8
	CUCTimeSeconds  uint32
	CUCTimeFraction uint16
}

// Read TMDataFieldHeader
func (tmdfh *TMDataFieldHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, tmdfh)
}

// Time returns the telemetry data time in UTC
func (tmdfh *TMDataFieldHeader) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(tmdfh.CUCTimeSeconds, tmdfh.CUCTimeFraction, epoch)
}

// Nanoseconds returns the telemetry data time in nanoseconds since its epoch
func (tmdfh *TMDataFieldHeader) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(tmdfh.CUCTimeSeconds, tmdfh.CUCTimeFraction)
}

// IsHousekeeping returns true if payload contains housekeeping data
func (tmdfh *TMDataFieldHeader) IsHousekeeping() bool {
	return tmdfh.ServiceType == HousekeepingDiagnosticDataReporting && tmdfh.ServiceSubType == 25
}

// IsTransparentData can be either CCD or Photometer data
func (tmdfh *TMDataFieldHeader) IsTransparentData() bool {
	return tmdfh.ServiceType == 128 && tmdfh.ServiceSubType == 25
}

// IsTCVerification returns true if payload contains TC verification data
func (tmdfh *TMDataFieldHeader) IsTCVerification() bool {
	return tmdfh.ServiceType == TelecommandVerification && (tmdfh.ServiceSubType == 1 ||
		tmdfh.ServiceSubType == 2 ||
		tmdfh.ServiceSubType == 7 ||
		tmdfh.ServiceSubType == 8)
}

// CSVHeaders returns the header row
func (tmdfh TMDataFieldHeader) CSVHeaders() []string {
	return []string{
		"TMHeaderTime",
		"TMHeaderNanoseconds",
	}
}

// CSVRow returns the data row
func (tmdfh TMDataFieldHeader) CSVRow() []string {
	gpsTime := time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
	tmTime := tmdfh.Time(gpsTime)
	return []string{
		tmTime.Format(time.RFC3339Nano),
		fmt.Sprintf("%v", tmdfh.Nanoseconds()),
	}
}
