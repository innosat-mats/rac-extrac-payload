package main

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

func readTCDFH(buf io.Reader) TCDataFieldHeader {
	tcdfh := TCDataFieldHeader{}
	err := binary.Read(buf, binary.BigEndian, &tcdfh)
	if err != nil {
		log.Fatal(err)
	}
	return tcdfh
}
func (dh TCDataFieldHeader) getPUS() uint8 {
	return (dh.PUS << 1) >> 5
}

func (dh TCDataFieldHeader) getServiceType() uint8 {
	return dh.ServiceType
}

func (dh TCDataFieldHeader) getServiceSubType() uint8 {
	return dh.ServiceSubType
}
