package ramses

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// Ramses data header
type Ramses struct {
	Synch  uint16 // fixed to 0xEB90
	Length uint16 // payload length
	Port   uint16
	Type   uint8
	Secure uint8
	Time   uint32 // milliseconds since midnight
	Date   int32  // days since 2000-01-01 00:00:00.00
}

// Read Ramses reads buffer into the struct
//
// returns io.EOF is there was exactly 0 bytes to read since this
// indicates previous record was the last record of the buffer and
// thus not really an error.
//
// for all other incomplete reads non EOF errors are produced.
func (ramses *Ramses) Read(buf io.Reader) error {
	size := binary.Size(ramses)
	tmpBuf := make([]byte, size)
	n, err := buf.Read(tmpBuf)
	if err != nil && err != io.EOF {
		return err
	}
	if n == 0 {
		return io.EOF
	} else if n != size {
		return errors.New("not enough data to read Ramses header")
	}

	return binary.Read(bytes.NewReader(tmpBuf), binary.LittleEndian, ramses)
}

// Created is when the package was created
func (ramses *Ramses) Created() time.Time {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return start.Add(time.Duration(ramses.Nanoseconds()))
}

// Nanoseconds since epoch when packet was created
func (ramses *Ramses) Nanoseconds() int64 {
	return int64(
		time.Hour*24*time.Duration(ramses.Date) +
			time.Millisecond*time.Duration(ramses.Time),
	)

}

// Valid as in correct version
func (ramses *Ramses) Valid() bool {
	return ramses.Synch == 0xeb90
}

// CSVSpecifications returns the specs used in creating the struct
func (ramses Ramses) CSVSpecifications() []string {
	return []string{"RAMSES", Specification}
}

//CSVHeaders returns the field names
func (ramses Ramses) CSVHeaders() []string {
	return []string{
		"RamsesTime",
	}
}

//CSVRow returns the field values
func (ramses Ramses) CSVRow() []string {
	return []string{
		fmt.Sprintf("%v", ramses.Created().Format(time.RFC3339Nano)),
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (ramses *Ramses) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Specification string `json:"specification"`
		RamsesTime    string `json:"ramsesTime"`
	}{
		Specification: Specification,
		RamsesTime:    ramses.Created().Format(time.RFC3339Nano),
	})
}
