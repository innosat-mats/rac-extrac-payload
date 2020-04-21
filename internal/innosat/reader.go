package innosat

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/howeyc/crc16"
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

// DecodeSource decodes a byte array to a SourcePackage
func DecodeSource(data []byte) (SourcePackage, error) {
	var err error
	buf := bytes.NewReader(data)
	header := SourcePacketHeader{}
	err = header.Read(buf)
	if err != nil {
		return SourcePackage{}, err
	}
	if crc16.ChecksumCCITTFalse(data[:len(data)-2]) != binary.BigEndian.Uint16(data[len(data)-2:]) {
		return SourcePackage{}, errors.New("checksum bad")
	}
	var payload SourcePayload
	if header.Type() == TM {
		tmpayload := TMDataFieldHeader{}
		err = tmpayload.Read(buf)
		if err != nil {
			return SourcePackage{}, err
		}
		payload = tmpayload
	}
	if header.Type() == TC {
		tcpayload := TCDataFieldHeader{}
		err = tcpayload.Read(buf)
		if err != nil {
			return SourcePackage{}, err
		}
		payload = tcpayload
	}

	return SourcePackage{
		header,
		payload,
		data[binary.Size(header)+binary.Size(payload) : len(data)-2],
	}, nil
}
