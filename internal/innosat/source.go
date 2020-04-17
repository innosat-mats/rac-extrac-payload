package innosat

import (
	"encoding/binary"
	"io"
)

//SourcePacketHeaderType is the type of the source packet (TM/TC)
type SourcePacketHeaderType uint

//TM is the source package type for telemetry
//
//TC is the source package type for telecommand
const (
	TC SourcePacketHeaderType = 0
	TM SourcePacketHeaderType = 1
)

//SourcePacketHeader Source Packet Header
type SourcePacketHeader struct {
	PacketID              uint16
	PacketSequenceControl uint16
	PacketLength          uint16
}

//Read SourcePacketHeader
func (sph *SourcePacketHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, sph)
}

// Version ...
func (sph *SourcePacketHeader) Version() uint {
	return uint(sph.PacketID >> 13)
}

// Type is either Telecommand or Telemetry
func (sph *SourcePacketHeader) Type() SourcePacketHeaderType {
	return SourcePacketHeaderType((sph.PacketID << 3) >> 15)
}

// HeaderType ...
func (sph *SourcePacketHeader) HeaderType() uint {
	return uint((sph.PacketID << 4) >> 15)
}

func (sph *SourcePacketHeader) apid() uint16 {
	return (sph.PacketID << 5) >> 5
}

// IsMainApplication says if packet is for main application
func (sph *SourcePacketHeader) IsMainApplication() bool {
	return sph.apid() == 100
}

// GroupingFlags ...
func (sph *SourcePacketHeader) GroupingFlags() uint {
	return uint(sph.PacketSequenceControl >> 14)
}

// SequenceCount ...
func (sph *SourcePacketHeader) SequenceCount() uint16 {
	return (sph.PacketSequenceControl << 2) >> 2
}
