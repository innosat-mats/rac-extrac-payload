package aez

import (
	"encoding/json"
	"fmt"
	"time"
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

func (sid *SID) String() string {
	switch *sid {
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
		return fmt.Sprintf("Unknown SID: %v", int(*sid))
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

// IsCCD returns if the RID is for a CCD
func (rid *RID) IsCCD() bool {
	return *rid == CCD1 || *rid == CCD2 || *rid == CCD3 || *rid == CCD4 || *rid == CCD5 || *rid == CCD6 || *rid == CCD7
}

func (rid *RID) String() string {
	switch *rid {
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
		return fmt.Sprintf("Unknown RID: %v", int(*rid))
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (rid *RID) MarshalJSON() ([]byte, error) {
	return json.Marshal(rid.String())
}
