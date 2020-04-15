package ramses

import (
	"encoding/binary"
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

// Read Ramses
func (r *Ramses) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, r)
}

// Created is when the package was created
func (r *Ramses) Created() time.Time {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return start.Add(time.Hour * 24 * time.Duration(r.Date)).Add(time.Millisecond * time.Duration(r.Time))
}

// Valid as in correct version
func (r *Ramses) Valid() bool {
	return r.Synch == 0xeb90
}

//SecureTrans always true?
func (r *Ramses) SecureTrans() bool {
	return true
}
