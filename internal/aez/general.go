package aez

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

// Specification describes what version the current implementation follows
var Specification string = "AEZICD002:H"

var gpsTime time.Time = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)

// SID is the id of a single housekeeping parameter
type SID uint16

const (
	// SIDSTAT is the SID of STAT.
	SIDSTAT SID = 1
	// SIDHTR is the SID of HTR.
	SIDHTR SID = 10
	// SIDPWR is the SID of PWR.
	SIDPWR SID = 20
	// SIDCPRUA is the SID of CPRUA.
	SIDCPRUA SID = 30
	// SIDCPRUB is the SID of CPRUB.
	SIDCPRUB SID = 31
)

func (sid SID) String() string {
	switch sid {
	case 0:
		return ""
	case SIDSTAT:
		return "STAT"
	case SIDHTR:
		return "HTR"
	case SIDPWR:
		return "PWR"
	case SIDCPRUA:
		return "CPRUA"
	case SIDCPRUB:
		return "CPRUB"
	default:
		return fmt.Sprintf("Unknown SID: %v", int(sid))
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (sid *SID) MarshalJSON() ([]byte, error) {
	return json.Marshal(sid.String())
}

// RID is Report Identification
type RID uint16

const (
	// CCD1 is connected to CPRUA port 0
	CCD1 RID = 21
	// CCD2 is connected to CPRUA port 1
	CCD2 RID = 22
	// CCD3 is connected to CPRUA port 2
	CCD3 RID = 23
	// CCD4 is connected to CPRUA port 3
	CCD4 RID = 24
	// CCD5 is connected to CPRUB port 0
	CCD5 RID = 25
	// CCD6 is connected to CPRUB port 1
	CCD6 RID = 26
	// CCD7 is connected to CPRUB port 2
	CCD7 RID = 27
	// PM is Photometer data
	PM RID = 30
)

func (rid RID) String() string {
	switch rid {
	case 0:
		return ""
	case CCD1:
		return "CCD1"
	case CCD2:
		return "CCD2"
	case CCD3:
		return "CCD3"
	case CCD4:
		return "CCD4"
	case CCD5:
		return "CCD5"
	case CCD6:
		return "CCD6"
	case CCD7:
		return "CCD7"
	case PM:
		return "PM"
	default:
		return fmt.Sprintf("Unknown RID: %v", int(rid))
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (rid *RID) MarshalJSON() ([]byte, error) {
	return json.Marshal(rid.String())
}

//STAT General status housekeeping report of the payload instrument.
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

// Read STAT
func (stat *STAT) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, stat)
}

// Time returns the measurement time in UTC
func (stat *STAT) Time(epoch time.Time) time.Time {
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
func (stat STAT) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "STATTIME", "STATNANO")
	// We don't need the raw CUC Time fields, instead the iso date and nanoseconds are included above.
	return append(headers, csvHeader(stat, "TS", "TSS")...)
}

// CSVRow returns the data row
func (stat STAT) CSVRow() []string {
	var row []string
	statTime := stat.Time(gpsTime)
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
