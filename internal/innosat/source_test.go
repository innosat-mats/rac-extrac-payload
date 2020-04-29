package innosat

import (
	"testing"
)

func TestVersion(t *testing.T) {
	sph := SourcePacketHeader{
		0b1110000000000000,
		0,
		0}

	ans := sph.PacketID.Version()
	if ans != 0b111 {
		t.Errorf("SourcePacketHeader.PacketID.Version() = %03b; want 111", ans)
	}
}

func TestType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0001000000000000,
		0,
		0}

	ans := sph.PacketID.Type()
	if ans != TC {
		t.Errorf("SourcePacketHeader.PacketID.Type() = %b; want 1", ans)
	}
}
func TestHeaderType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000100000000000,
		0,
		0}

	ans := sph.PacketID.HeaderType()
	if ans != 1 {
		t.Errorf("SourcePacketHeader.PacketID.HeaderType() = %b; want 1", ans)
	}
}

func TestGroupingFlags(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b1100000000000000,
		0}

	ans := sph.PacketSequenceControl.GroupingFlags()
	if ans != 0b11 {
		t.Errorf("SourcePacketHeader.PacketSequenceControl.GroupingFlags() = %02b; want 11", ans)
	}
}

func TestSequenceCount(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b0011111111111111,
		0}

	ans := sph.PacketSequenceControl.SequenceCount()
	if ans != 0b11111111111111 {
		t.Errorf("SourcePacketHeader.PacketSequenceControl.SequenceCount() = %02b; want 11111111111111", ans)
	}
}

func TestSourcePacketHeader_APID(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   SourcePacketAPIDType
	}{
		{"Returns IdleAPID when apid is 0", fields{}, TimeAPID},
		{"Returns MainAPID when apid is 100", fields{PacketID: 2148}, MainAPID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			if got := sph.PacketID.APID(); got != tt.want {
				t.Errorf("SourcePacketHeader.PackedID.APID() = %v, want %v", got, tt.want)
			}
		})
	}
}
