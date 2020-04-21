package innosat

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/howeyc/crc16"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// SourcePayload ...
type SourcePayload interface {
	PUSVersion() uint8
}

// SourcePackage ...
type SourcePackage struct {
	Header      SourcePacketHeader
	Payload     SourcePayload
	Application []byte
}

// DecodeSource ..
func DecodeSource(ramsesPackage ramses.Package) SourcePackage {
	var err error
	buf := bytes.NewReader(ramsesPackage.Payload)
	header := SourcePacketHeader{}
	err = header.Read(buf)
	if err != nil {
		log.Fatalln("Source Header read error", err)
	}
	if crc16.ChecksumCCITTFalse(ramsesPackage.Payload[:len(ramsesPackage.Payload)-2]) != binary.BigEndian.Uint16(ramsesPackage.Payload[len(ramsesPackage.Payload)-2:]) {
		log.Fatal("checksum bad")
	}
	var payload SourcePayload
	if header.Type() == TM {
		tmpayload := TMDataFieldHeader{}
		err = tmpayload.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		payload = tmpayload
	}
	if header.Type() == TC {
		tcpayload := TCDataFieldHeader{}
		err = tcpayload.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		payload = tcpayload
	}

	return SourcePackage{
		header,
		payload,
		ramsesPackage.Payload[binary.Size(header)+binary.Size(payload) : len(ramsesPackage.Payload)-2],
	}
}
