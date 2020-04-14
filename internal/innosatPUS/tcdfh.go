package pus

import (
	"encoding/binary"
	"io"
	"log"
)

//TCDataFieldHeader (9 octets)
type TCDataFieldHeader struct {
	PUS            uint8
	ServiceType    uint8
	ServiceSubType uint8
}

// TCDataFieldHeader
func readTCDFH(buf io.Reader) TCDataFieldHeader {
	tcdfh := TCDataFieldHeader{}
	err := binary.Read(buf, binary.BigEndian, &tcdfh)
	if err != nil {
		log.Fatal(err)
	}
	return tcdfh
}

//PUSVersion ...
func (dh TCDataFieldHeader) PUSVersion() uint8 {
	return (dh.PUS << 1) >> 5
}
