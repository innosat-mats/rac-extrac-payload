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
	Header             *innosat.SourcePacketHeader
	Payload            *innosat.TMHeader
	ApplicationPayload []byte
}

const crcChecksumLength int = 2
const pusLengthOffset int = 1

// DecodeSource decodes a byte array to a SourcePackage
func DecodeSource(data []byte) (SourcePackage, error) {
	var err error
	buf := bytes.NewReader(data)
	header, err := innosat.NewSourcePacketHeader(buf)
	if err != nil {
		return SourcePackage{}, err
	}
	if crc16.ChecksumCCITTFalse(data[:len(data)-crcChecksumLength]) != binary.BigEndian.Uint16(data[len(data)-crcChecksumLength:]) {
		return SourcePackage{}, fmt.Errorf(
			"checksum bad %d",
			crc16.ChecksumCCITTFalse(data[:len(data)-crcChecksumLength]),
		)
	}

	tmHeader, err := innosat.NewTMHeader(buf)
	if err != nil {
		return SourcePackage{}, err
	}

	sliceStart := binary.Size(header) + binary.Size(tmHeader)
	sliceEnd := sliceStart + int(header.PacketLength) - binary.Size(tmHeader) - crcChecksumLength + pusLengthOffset
	return SourcePackage{
		header,
		tmHeader,
		data[sliceStart:sliceEnd],
	}, nil
}
