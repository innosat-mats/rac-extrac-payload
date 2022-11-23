package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
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
// SetParquet sets the parquet representation of the STAT
func (stat *STAT) SetParquet(row *parquetrow.ParquetRow) {
	row.STATTime = stat.Time(GpsTime)
	row.STATNanoseconds = stat.Nanoseconds()
	row.SPID = stat.SPID
	row.SPREV = stat.SPREV
	row.FPID = stat.FPID
	row.FPREV = stat.FPREV
	row.SVNA = stat.SVNA
	row.SVNB = stat.SVNB
	row.SVNC = stat.SVNC
	row.MODE = stat.MODE
	row.EDACE = stat.EDACE
	row.EDACCE = stat.EDACCE
	row.EDACN = stat.EDACN
	row.SPWEOP = stat.SPWEOP
	row.SPWEEP = stat.SPWEEP
	row.ANOMALY = stat.ANOMALY
}
