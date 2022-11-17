package innosat

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// SourcePacketHeaderType is the type of the source packet (TM/TC)
type SourcePacketHeaderType uint

const (
	//TM is the source package type for telemetry
	TM SourcePacketHeaderType = 0
	//TC is the source package type for telecommand
	TC SourcePacketHeaderType = 1
)

func (headerType *SourcePacketHeaderType) String() string {
	switch *headerType {
	case TM:
		return "TM"
	case TC:
		return "TC"
	default:
		return fmt.Sprintf("Unknown HeaderType %v", uint(*headerType))
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

func (continuationFlag *SourcePackageContinuationFlagType) String() string {
	switch *continuationFlag {
	case SPCont:
		return "Continuation"
	case SPStart:
		return "Start"
	case SPStop:
		return "Stop"
	case SPStandalone:
		return "Standalone"
	default:
		return fmt.Sprintf("Unknown ContinuationFlag %v", uint(*continuationFlag))
	}
}

type packetID uint16

// Version ...
func (pid *packetID) Version() uint {
	return uint(*pid >> 13)
}

// Type is either Telecommand or Telemetry
func (pid *packetID) Type() SourcePacketHeaderType {
	return SourcePacketHeaderType((*pid << 3) >> 15)
}

// HeaderType ...
func (pid *packetID) HeaderType() uint {
	return uint((*pid << 4) >> 15)
}

// APID ...
func (pid *packetID) APID() SourcePacketAPIDType {
	return SourcePacketAPIDType(*pid & 0x07FF)
}

// PacketSequenceControl is the encoding of the sequence value
type PacketSequenceControl uint16

// GroupingFlags ...
func (psc *PacketSequenceControl) GroupingFlags() SourcePackageContinuationFlagType {
	return SourcePackageContinuationFlagType(*psc >> 14)
}

// SequenceCount ...
func (psc *PacketSequenceControl) SequenceCount() uint16 {
	return uint16((*psc << 2) >> 2)
}

// SourcePacketHeader Source Packet Header
type SourcePacketHeader struct {
	PacketID              packetID
	PacketSequenceControl PacketSequenceControl
	PacketLength          uint16
}

// NewSourcePacketHeader reads a SourcePacketHeader from buffer
func NewSourcePacketHeader(buf io.Reader) (*SourcePacketHeader, error) {
	sph := SourcePacketHeader{}
	err := binary.Read(buf, binary.BigEndian, &sph)
	return &sph, err
}

// CSVSpecifications returns the version of the spec used
func (sph *SourcePacketHeader) CSVSpecifications() []string {
	return []string{"INNOSAT", Specification}
}

// CSVHeaders returns the header row
func (sph *SourcePacketHeader) CSVHeaders() []string {
	return []string{
		"SPSequenceCount",
	}
}

// CSVRow returns the data row
func (sph *SourcePacketHeader) CSVRow() []string {
	return []string{
		fmt.Sprintf("%v", sph.PacketSequenceControl.SequenceCount()),
	}
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (sph *SourcePacketHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Specification   string `json:"specification"`
		SPSequenceCount uint16 `json:"spSequenceCount"`
	}{
		Specification:   Specification,
		SPSequenceCount: sph.PacketSequenceControl.SequenceCount(),
	})
}

// SourcePackedHeaderParquet holds the parquet representation of the SourcePacketHeader
type SourcePacketHeaderParquet struct {
	SPSequenceCount uint16 `parquet:"SPSequenceCount"`
}

// GetParquet returns the parquet representation of the SourcePacketHeader
func (sph *SourcePacketHeader) GetParquet() SourcePacketHeaderParquet {
	return SourcePacketHeaderParquet{
		sph.PacketSequenceControl.SequenceCount(),
	}
}
