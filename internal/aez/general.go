package aez

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)


// Specification describes what version the current implementation follows
var Specification string = "AEZICD002:E"

// SID is the id of a single housekeeping parameter
type SID uint16

// SIDSTAT is the SID of STAT.
//
// SIDHTR is the SID of HTR.
//
// SIDPWR is the SID of PWR.
//
// SIDCPRUA is the SID of CPRUA.
//
// SIDCPRUB is the SID of CPRUB.
const (
	SIDSTAT  SID = 1
	SIDHTR   SID = 10
	SIDPWR   SID = 20
	SIDCPRUA SID = 30
	SIDCPRUB SID = 31
)

//STAT General status housekeeping report of the payload instrument.
type STAT struct { //(34 octets)
	SPID   uint16 // Software Part ID
	SPREV  uint8  // Software Part Revision
	FPID   uint16 // Firmware Part ID
	FPREV  uint8  // Firmware Part Revision
	TS     uint32 // Time, seconds (CUC time format)
	TSS    uint16 // Time, subseconds (CUC time format)
	MODE   uint8  // Payload mode 1..2
	EDACE  uint32 // EDAC detected single bit errors
	EDACCE uint32 // EDAC corrected single bit errors
	EDACN  uint32 // EDAC memory scrubber passes through memory
	SPWEOP uint32 // SpaceWire received EOPs
	SPWEEP uint32 // SpaceWire received EEPs
}

// Read STAT
func (stat *STAT) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, stat)
}

// Time returns the measurement time in UTC
func (stat *STAT) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(stat.TS, stat.TSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (stat *STAT) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(stat.TS, stat.TSS)
}
