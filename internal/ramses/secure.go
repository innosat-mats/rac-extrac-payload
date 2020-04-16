package ramses

import (
	"encoding/binary"
	"io"
)

// Secure header
type Secure struct {
	IPAddress      uint32
	Port           uint16
	Seq            uint16
	Retransmission uint16
	Ack            uint16
	_              uint32
}

// Read fills a structure with data
func (secure *Secure) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, secure)
}
