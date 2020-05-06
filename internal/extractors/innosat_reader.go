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

// DecodeSource decodes a byte array to a SourcePackage
func DecodeSource(data []byte) (SourcePackage, error) {
	var err error
	buf := bytes.NewReader(data)
	header := innosat.SourcePacketHeader{}
	err = header.Read(buf)
	if err != nil {
		return SourcePackage{}, err
	}
	if crc16.ChecksumCCITTFalse(data[:len(data)-2]) != binary.BigEndian.Uint16(data[len(data)-2:]) {
		return SourcePackage{}, fmt.Errorf("checksum bad %d", crc16.ChecksumCCITTFalse(data[:len(data)-2]))
	}

	tmPayload := innosat.TMDataFieldHeader{}
	err = tmPayload.Read(buf)
	if err != nil {
		return SourcePackage{}, err
	}

	return SourcePackage{
		header,
		tmPayload,
		data[binary.Size(header)+binary.Size(tmPayload) : len(data)-2],
	}, nil
}
