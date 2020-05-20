package timeseries

import (
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// OutStream is the type for the outstream enum
type OutStream int

const (
	// Unknown is an unregonized out stream
	Unknown OutStream = iota
	// HTR is a HTR timeseries out stream
	HTR
	// PWR is a PWR timeseries out stream
	PWR
	// CPRU is a CPRU timeseries out stream
	CPRU
	// STAT is as STAT timeseries out stream
	STAT
	// PM is a PM timeseries out stream
	PM
	// CCD is a CCD timeseries out stream
	CCD
	// TCV is a TCV timeseries out stream
	TCV
)

func (stream OutStream) String() string {
	switch stream {
	case HTR:
		return "HTR"
	case PWR:
		return "PWR"
	case CPRU:
		return "CPRU"
	case STAT:
		return "STAT"
	case PM:
		return "PM"
	case CCD:
		return "CCD"
	case TCV:
		return "TCV"
	default:
		return "unknown"
	}
}

// OutStreamFromDataRecord infers stream based on data
func OutStreamFromDataRecord(pkg *common.DataRecord) OutStream {
	switch pkg.Data.(type) {
	case *aez.CCDImage:
		return CCD
	case *aez.PMData:
		return PM
	case *aez.HTR:
		return HTR
	case *aez.PWR:
		return PWR
	case *aez.CPRU:
		return CPRU
	case *aez.STAT:
		return STAT
	case *aez.TCAcceptSuccessData, *aez.TCAcceptFailureData, *aez.TCExecSuccessData, *aez.TCExecFailureData:
		return TCV
	default:
		return Unknown
	}
}
