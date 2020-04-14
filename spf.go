package main

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

func (h *SourcePacketHeader) read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, h)
}

func (h *SourcePacketHeader) getVersion() uint {
	return uint(h.PacketID >> 13)
}

func (h *SourcePacketHeader) getType() uint {
	return uint((h.PacketID << 3) >> 15)
}

func (h *SourcePacketHeader) getHeaderType() uint {
	return uint((h.PacketID << 4) >> 15)
}

func (h *SourcePacketHeader) getAPID() uint16 {
	return (h.PacketID << 5) >> 5
}

func (h *SourcePacketHeader) getGroupingFlags() uint {
	return uint(h.PacketSequenceControl >> 14)
}

func (h *SourcePacketHeader) getSequenceCount() uint16 {
	return (h.PacketSequenceControl << 2) >> 2
}
