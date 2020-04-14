package pus

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

//Read ...
func (h *SourcePacketHeader) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, h)
}

// Version ...
func (h *SourcePacketHeader) Version() uint {
	return uint(h.PacketID >> 13)
}

// Type ...
func (h *SourcePacketHeader) Type() uint {
	return uint((h.PacketID << 3) >> 15)
}

// HeaderType ...
func (h *SourcePacketHeader) HeaderType() uint {
	return uint((h.PacketID << 4) >> 15)
}

// APID ...
func (h *SourcePacketHeader) APID() uint16 {
	return (h.PacketID << 5) >> 5
}

// GroupingFlags ...
func (h *SourcePacketHeader) GroupingFlags() uint {
	return uint(h.PacketSequenceControl >> 14)
}

// SequenceCount ...
func (h *SourcePacketHeader) SequenceCount() uint16 {
	return (h.PacketSequenceControl << 2) >> 2
}
