package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

// STAT General status housekeeping report of the payload instrument.
type STAT struct { //(34 octets)
	SPID    uint16 // Software Part ID
	SPREV   uint8  // Software Part Revision
	FPID    uint16 // Firmware Part ID
	FPREV   uint8  // Firmware Part Revision
	SVNA    uint8  // SVN revision tag on format A.B.C
	SVNB    uint8  // SVN revision tag on format A.B.C
	SVNC    uint8  // SVN revision tag on format A.B.C
	TS      uint32 // Time, seconds (CUC time format)
	TSS     uint16 // Time, subseconds (CUC time format)
	MODE    uint8  // Payload mode 1..2
	EDACE   uint32 // EDAC detected single bit errors
	EDACCE  uint32 // EDAC corrected single bit errors
	EDACN   uint32 // EDAC memory scrubber passes through memory
	SPWEOP  uint32 // SpaceWire received EOPs
	SPWEEP  uint32 // SpaceWire received EEPs
	ANOMALY uint8  // Anomalyflag (0==0 ? OK: payload power off)
}

// NewSTAT reads a STAT from buffer
func NewSTAT(buf io.Reader) (*STAT, error) {
	stat := STAT{}
	err := binary.Read(buf, binary.LittleEndian, &stat)
	return &stat, err
}

// Time returns the measurement time in UTC
func (stat *STAT) Time(epoch time.Time) time.Time {
	if (epoch == time.Time{}) {
		epoch = GpsTime
	}
	return ccsds.UnsegmentedTimeDate(stat.TS, stat.TSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (stat *STAT) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(stat.TS, stat.TSS)
}

// CSVSpecifications returns the version of the spec used
func (stat STAT) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (stat *STAT) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "STATTIME", "STATNANO")
	// We don't need the raw CUC Time fields, instead the iso date and nanoseconds are included above.
	return append(headers, csvHeader(stat, "TS", "TSS")...)
}

// CSVRow returns the data row
func (stat *STAT) CSVRow() []string {
	var row []string
	statTime := stat.Time(GpsTime)
	row = append(row, statTime.Format(time.RFC3339Nano), fmt.Sprintf("%v", stat.Nanoseconds()))
	val := reflect.Indirect(reflect.ValueOf(stat))
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		name := t.Field(i).Name
		// We don't need the raw CUC Time fields, instead the iso date and nanoseconds are included above.
		if name != "TS" && name != "TSS" {
			valueField := val.Field(i)
			row = append(row, fmt.Sprintf("%v", valueField.Uint()))
		}
	}
	return row
}

// STATParquet holds the parquet representation of the STAT
type STATParquet struct { //(34 octets)
	STATTime        time.Time `parquet:"STATTime"`
	STATNanoseconds int64     `parquet:"STATNanoseconds"`
	SPID            uint16    `parquet:"SPID"`
	SPREV           uint8     `parquet:"SPREV"`
	FPID            uint16    `parquet:"FPID"`
	FPREV           uint8     `parquet:"FPREV"`
	SVNA            uint8     `parquet:"SVNA"`
	SVNB            uint8     `parquet:"SVNB"`
	SVNC            uint8     `parquet:"SVNC"`
	MODE            uint8     `parquet:"MODE"`
	EDACE           uint32    `parquet:"EDACE"`
	EDACCE          uint32    `parquet:"EDACCE"`
	EDACN           uint32    `parquet:"EDACN"`
	SPWEOP          uint32    `parquet:"SPWEOP"`
	SPWEEP          uint32    `parquet:"SPWEEP"`
	ANOMALY         uint8     `parquet:"ANOMALY"`
}

// GetParquet returns the parquet representation of the STAT
func (stat *STAT) GetParquet() STATParquet {
	return STATParquet{
		stat.Time(GpsTime),
		stat.Nanoseconds(),
		stat.SPID,
		stat.SPREV,
		stat.FPID,
		stat.FPREV,
		stat.SVNA,
		stat.SVNB,
		stat.SVNC,
		stat.MODE,
		stat.EDACE,
		stat.EDACCE,
		stat.EDACN,
		stat.SPWEOP,
		stat.SPWEEP,
		stat.ANOMALY,
	}
}
