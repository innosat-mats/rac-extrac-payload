package innosat

import (
	"encoding/binary"
	"fmt"
	"io"
)

//SourcePacketHeaderType is the type of the source packet (TM/TC)
type SourcePacketHeaderType uint

const (
	//TM is the source package type for telemetry
	TM SourcePacketHeaderType = 0
	//TC is the source package type for telecommand
	TC SourcePacketHeaderType = 1
)

func (headerType SourcePacketHeaderType) String() string {
	switch headerType {
	case TM:
		return "TM"
	case TC:
		return "TC"
	default:
		return fmt.Sprintf("Unknown %v", headerType)
	}
}

// SourcePackageContinuationFlagType type for continuation groups
type SourcePackageContinuationFlagType uint

const (
	// SPCont Continuation packet
	SPCont SourcePackageContinuationFlagType = iota
	// SPStart start sequence of continuation packets
	SPStart
	// SPStop end of continuation packets
	SPStop
	// SPStandalone a single packet
	SPStandalone
)

func (continuationFlag SourcePackageContinuationFlagType) String() string {
	switch continuationFlag {
	case SPCont:
		return "Continuation"
	case SPStart:
		return "Start"
	case SPStop:
		return "Stop"
	case SPStandalone:
		return "Standalone"
	default:
		return fmt.Sprintf("Unknown %v", continuationFlag)
	}
}

type packetID uint16

// Version ...
func (pid packetID) Version() uint {
	return uint(pid >> 13)
}

// Type is either Telecommand or Telemetry
func (pid packetID) Type() SourcePacketHeaderType {
	return SourcePacketHeaderType((pid << 3) >> 15)
}

// HeaderType ...
func (pid packetID) HeaderType() uint {
	return uint((pid << 4) >> 15)
}

// APID ...
func (pid packetID) APID() SourcePacketAPIDType {
	return SourcePacketAPIDType(pid & 0x07FF)
}

type packetSequenceControl uint16

// GroupingFlags ...
func (psc packetSequenceControl) GroupingFlags() SourcePackageContinuationFlagType {
	return SourcePackageContinuationFlagType(psc >> 14)
}

// SequenceCount ...
func (psc packetSequenceControl) SequenceCount() uint16 {
	return uint16((psc << 2) >> 2)
}

// SourcePacketHeader Source Packet Header
type SourcePacketHeader struct {
	PacketID              packetID
	PacketSequenceControl packetSequenceControl
	PacketLength          uint16
}

//Read SourcePacketHeader
func (sph *SourcePacketHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, sph)
}

// CSVSpecifications returns the version of the spec used
func (sph SourcePacketHeader) CSVSpecifications() []string {
	return []string{"INNOSAT", Specification}
}

// CSVHeaders returns the header row
func (sph SourcePacketHeader) CSVHeaders() []string {
	return []string{
		"SourcePacketVersion",
		"SourcePacketType",
		"SourcePacketHeaderType",
		"SourcePacketAPID",
		"SourcePacketGroupingFlags",
		"SourcePacketSequenceCount",
		"SourcePacketLength",
	}
}

// CSVRow returns the data row
func (sph SourcePacketHeader) CSVRow() []string {
	return []string{
		fmt.Sprintf("%v", sph.PacketID.Version()),
		sph.PacketID.Type().String(),
		fmt.Sprintf("%v", sph.PacketID.HeaderType()),
		fmt.Sprintf("%v", sph.PacketID.APID()),
		sph.PacketSequenceControl.GroupingFlags().String(),
		fmt.Sprintf("%v", sph.PacketSequenceControl.SequenceCount()),
		fmt.Sprintf("%v", sph.PacketLength),
	}
}
