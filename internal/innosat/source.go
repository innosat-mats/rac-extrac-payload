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
	TM SourcePacketHeaderType = 0
	TC SourcePacketHeaderType = 1
)

//SourcePackageContinuationFlagType type for continuation groups
type SourcePackageContinuationFlagType uint

//SPCont Continuation packet
//SPStart start sequence of continuation packets
//SPStop end of continuation packets
//SPStandalone a single packet
const (
	SPCont SourcePackageContinuationFlagType = iota
	SPStart
	SPStop
	SPStandalone
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
func (sph SourcePacketHeader) Type() SourcePacketHeaderType {
	return SourcePacketHeaderType((sph.PacketID << 3) >> 15)
}

// HeaderType ...
func (sph *SourcePacketHeader) HeaderType() uint {
	return uint((sph.PacketID << 4) >> 15)
}

// APID ...
func (sph SourcePacketHeader) APID() uint16 {
	return (sph.PacketID << 5) >> 5
}

// IsMainApplication says if packet is for main application
func (sph *SourcePacketHeader) IsMainApplication() bool {
	return sph.APID() == 100
}

// GroupingFlags ...
func (sph SourcePacketHeader) GroupingFlags() SourcePackageContinuationFlagType {
	return SourcePackageContinuationFlagType(sph.PacketSequenceControl >> 14)
}

// SequenceCount ...
func (sph *SourcePacketHeader) SequenceCount() uint16 {
	return (sph.PacketSequenceControl << 2) >> 2
}
