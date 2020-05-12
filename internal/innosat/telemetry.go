package innosat

import (
	"encoding/binary"
	"encoding/json"
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

//TMHeader (9 octets)
type TMHeader struct {
	PUS             pus
	ServiceType     SourcePackageServiceType
	ServiceSubType  uint8
	CUCTimeSeconds  uint32
	CUCTimeFraction uint16
}

// Read TMHeader
func (header *TMHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, header)
}

// Time returns the telemetry data time in UTC
func (header *TMHeader) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(header.CUCTimeSeconds, header.CUCTimeFraction, epoch)
}

// Nanoseconds returns the telemetry data time in nanoseconds since its epoch
func (header *TMHeader) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(header.CUCTimeSeconds, header.CUCTimeFraction)
}

// IsHousekeeping returns true if payload contains housekeeping data
func (header *TMHeader) IsHousekeeping() bool {
	return header.ServiceType == HousekeepingDiagnosticDataReporting && header.ServiceSubType == 25
}

// IsTransparentData can be either CCD or Photometer data
func (header *TMHeader) IsTransparentData() bool {
	return header.ServiceType == 128 && header.ServiceSubType == 25
}

// IsTCVerification returns true if payload contains TC verification data
func (header *TMHeader) IsTCVerification() bool {
	return header.ServiceType == TelecommandVerification &&
		(header.ServiceSubType == 1 ||
			header.ServiceSubType == 2 ||
			header.ServiceSubType == 7 ||
			header.ServiceSubType == 8)
}

// CSVHeaders returns the header row
func (header TMHeader) CSVHeaders() []string {
	return []string{
		"TMHeaderTime",
		"TMHeaderNanoseconds",
	}
}

var gpsTime = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)

// CSVRow returns the data row
func (header TMHeader) CSVRow() []string {
	tmTime := header.Time(gpsTime)
	return []string{
		tmTime.Format(time.RFC3339Nano),
		fmt.Sprintf("%v", header.Nanoseconds()),
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (header *TMHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TMHeaderTime        string `json:"tmHeaderTime"`
		TMHeaderNanoseconds int64  `json:"tmHeaderNanoseconds"`
	}{
		TMHeaderTime:        header.Time(gpsTime).Format(time.RFC3339Nano),
		TMHeaderNanoseconds: header.Nanoseconds(),
	})
}
