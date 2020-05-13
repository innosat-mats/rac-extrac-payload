package aez

import (
	"encoding/json"
	"fmt"
	"time"
)

// CCDImage is a container for the invariant CCDImagePackData header and the variable BadColumns that follow
type CCDImage struct {
	PackData   CCDImagePackData
	BadColumns []uint16
}

// CSVSpecifications returns the specs used in creating the struct
func (ccd CCDImage) CSVSpecifications() []string {
	return []string{"Specification", Specification}
}

// CSVHeaders returns the exportable field names
func (ccd CCDImage) CSVHeaders() []string {
	return append(ccd.PackData.CSVHeaders(), "BC")
}

// CSVRow returns the exportable field values
func (ccd CCDImage) CSVRow() []string {
	row := ccd.PackData.CSVRow()
	return append(row, fmt.Sprintf("%v", ccd.BadColumns))
}

// MarshalJSON jsonifies content
func (ccd *CCDImage) MarshalJSON() ([]byte, error) {
	wdwhigh, wdwlow, _ := ccd.PackData.WDW.InputDataWindow()
	return json.Marshal(&struct {
		Specification      string `json:"specification"`
		CCDSEL             uint8
		EXPNanoseconds     int64
		EXPDate            string
		WDWMode            string
		WDWInputDataWindow string
		WDWOV              uint16
		JPEGQ              uint8
		FRAME              uint16
		NROW               uint16
		NRBIN              uint16
		NRSKIP             uint16
		NCOL               uint16
		NCBINFPGAColumns   int
		NCBINCCDColumns    int
		NCSKIP             uint16
		NFLUSH             uint16
		TEXPMS             uint32
		GAINMode           string
		GAINTiming         string
		TEMP               uint16
		FBINOV             uint16
		LBLNK              uint16
		TBLNK              uint16
		ZERO               uint16
		TIMING1            uint16
		TIMING2            uint16
		VERSION            uint16
		TIMING3            uint16
		NBC                uint16
		BC                 []uint16
	}{
		Specification,
		ccd.PackData.CCDSEL,
		ccd.PackData.Nanoseconds(),
		ccd.PackData.Time(gpsTime).Format(time.RFC3339Nano),
		ccd.PackData.WDW.Mode().String(),
		fmt.Sprintf("%v..%v", wdwhigh, wdwlow),
		ccd.PackData.WDWOV,
		ccd.PackData.JPEGQ,
		ccd.PackData.FRAME,
		ccd.PackData.NROW,
		ccd.PackData.NRBIN,
		ccd.PackData.NRSKIP,
		ccd.PackData.NCOL,
		ccd.PackData.NCBIN.FPGAColumns(),
		ccd.PackData.NCBIN.CCDColumns(),
		ccd.PackData.NCSKIP,
		ccd.PackData.NFLUSH,
		ccd.PackData.TEXPMS,
		ccd.PackData.GAIN.Mode().String(),
		ccd.PackData.GAIN.Timing().String(),
		ccd.PackData.TEMP,
		ccd.PackData.FBINOV,
		ccd.PackData.LBLNK,
		ccd.PackData.TBLNK,
		ccd.PackData.ZERO,
		ccd.PackData.TIMING1,
		ccd.PackData.TIMING2,
		ccd.PackData.VERSION,
		ccd.PackData.TIMING3,
		ccd.PackData.NBC,
		ccd.BadColumns,
	})
}
