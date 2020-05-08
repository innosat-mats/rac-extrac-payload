package extractors

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/howeyc/crc16"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

// SourcePackage ...
type SourcePackage struct {
	Header             innosat.SourcePacketHeader
	Payload            innosat.TMDataFieldHeader
	ApplicationPayload []byte
}

const crcChecksumLength int = 2
var pusLengthOffset int = 1

// DecodeSource decodes a byte array to a SourcePackage
func DecodeSource(data []byte) (SourcePackage, error) {
	var err error
	buf := bytes.NewReader(data)
	header := innosat.SourcePacketHeader{}
	err = header.Read(buf)
	if err != nil {
		return SourcePackage{}, err
	}
	if crc16.ChecksumCCITTFalse(data[:len(data)-crcChecksumLength]) != binary.BigEndian.Uint16(data[len(data)-crcChecksumLength:]) {
		return SourcePackage{}, fmt.Errorf(
			"checksum bad %d",
			crc16.ChecksumCCITTFalse(data[:len(data)-crcChecksumLength]),
		)
	}

	tmPayload := innosat.TMDataFieldHeader{}
	err = tmPayload.Read(buf)
	if err != nil {
		return SourcePackage{}, err
	}

	sliceStart := binary.Size(header) + binary.Size(tmPayload)
	sliceEnd := sliceStart + int(header.PacketLength) - binary.Size(tmPayload) - crcChecksumLength + pusLengthOffset
	return SourcePackage{
		header,
		tmPayload,
		data[sliceStart:sliceEnd],
	}, nil
}
