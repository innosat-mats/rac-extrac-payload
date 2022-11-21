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

// PMData data from photometers
type PMData struct {
	EXPTS    uint32 // Exposure start time, seconds (CUC time format)
	EXPTSS   uint16 // Exposure start time, subseconds (CUC time format)
	PM1A     uint32 // Photometer 1, thermistor input A sum
	PM1ACNTR uint32 // Photometer 1, thermistor input A counter
	PM1B     uint32 // Photometer 1, thermistor input B sum
	PM1BCNTR uint32 // Photometer 1, thermistor input B counter
	PM1S     uint32 // Photometer 1, photo diode input SIG sum
	PM1SCNTR uint32 // Photometer 1, photo diode input SIG counter
	PM2A     uint32 // Photometer 2, thermistor input A sum
	PM2ACNTR uint32 // Photometer 2, thermistor input A counter
	PM2B     uint32 // Photometer 2, thermistor input B sum
	PM2BCNTR uint32 // Photometer 2, thermistor input B counter
	PM2S     uint32 // Photometer 2, photo diode input SIG sum
	PM2SCNTR uint32 // Photometer 2, photo diode input SIG counter
}

// NewPMData reads a PMData from reader
func NewPMData(buf io.Reader) (*PMData, error) {
	pm := PMData{}
	err := binary.Read(buf, binary.LittleEndian, &pm)
	return &pm, err
}

// Time returns the measurement time in UTC
func (pm *PMData) Time(epoch time.Time) time.Time {
	if (epoch == time.Time{}) {
		epoch = GpsTime
	}
	return ccsds.UnsegmentedTimeDate(pm.EXPTS, pm.EXPTSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (pm *PMData) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(pm.EXPTS, pm.EXPTSS)
}

// CSVSpecifications returns the version of the spec used
func (pm *PMData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (pm *PMData) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "PMTIME", "PMNANO")
	// We don't need the raw CUC Time fields, instead the iso date and
	// nanoseconds are included above.
	return append(headers, csvHeader(pm, "EXPTS", "EXPTSS")...)
}

// CSVRow returns the data row
func (pm *PMData) CSVRow() []string {
	var row []string
	pmTime := pm.Time(GpsTime)
	row = append(row, pmTime.Format(time.RFC3339Nano), fmt.Sprintf("%v", pm.Nanoseconds()))
	val := reflect.Indirect(reflect.ValueOf(pm))
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		name := t.Field(i).Name
		// We don't need the raw CUC Time fields, instead the iso date and
		// nanoseconds are included above.
		if name != "EXPTS" && name != "EXPTSS" {
			valueField := val.Field(i)
			row = append(row, fmt.Sprintf("%v", valueField.Uint()))
		}
	}
	return row
}

// SetParquet sets the parquet representation of the PMData
func (pm *PMData) SetParquet(row *parquetrow.ParquetRow) {
	row.PMTime = pm.Time(GpsTime)
	row.PMNanoseconds = pm.Nanoseconds()
	row.PM1A = pm.PM1A
	row.PM1ACNTR = pm.PM1ACNTR
	row.PM1B = pm.PM1B
	row.PM1BCNTR = pm.PM1BCNTR
	row.PM1S = pm.PM1S
	row.PM1SCNTR = pm.PM1SCNTR
	row.PM2A = pm.PM2A
	row.PM2ACNTR = pm.PM2ACNTR
	row.PM2B = pm.PM2B
	row.PM2BCNTR = pm.PM2BCNTR
	row.PM2S = pm.PM2S
	row.PM2SCNTR = pm.PM2SCNTR
}
