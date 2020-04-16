package innosat

import (
	"encoding/binary"
	"io"
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

// Type ...
func (sph *SourcePacketHeader) Type() uint {
	return uint((sph.PacketID << 3) >> 15)
}

// HeaderType ...
func (sph *SourcePacketHeader) HeaderType() uint {
	return uint((sph.PacketID << 4) >> 15)
}

// APID ...
func (sph *SourcePacketHeader) APID() uint16 {
	return (sph.PacketID << 5) >> 5
}

// GroupingFlags ...
func (sph *SourcePacketHeader) GroupingFlags() uint {
	return uint(sph.PacketSequenceControl >> 14)
}

// SequenceCount ...
func (sph *SourcePacketHeader) SequenceCount() uint16 {
	return (sph.PacketSequenceControl << 2) >> 2
}
