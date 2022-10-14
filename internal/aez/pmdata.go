package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
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
	const gpsTimeCorrection = -18 // Seconds
	var row []string
	gpsTime := time.Date(1980, time.January, 6, 0, 0, gpsTimeCorrection, 0, time.UTC)
	pmTime := pm.Time(gpsTime)
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
