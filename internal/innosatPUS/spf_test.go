package pus

import (
	"testing"
)

func TestVersion(t *testing.T) {
	sph := SourcePacketHeader{
		0b1110000000000000,
		0,
		0}

	ans := sph.Version()
	if ans != 0b111 {
		t.Errorf("Version() = %03b; want 111", ans)
	}
}

func TestType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0001000000000000,
		0,
		0}

	ans := sph.Type()
	if ans != 1 {
		t.Errorf("Type() = %b; want 1", ans)
	}
}
func TestHeaderType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000100000000000,
		0,
		0}

	ans := sph.HeaderType()
	if ans != 1 {
		t.Errorf("HeaderType() = %b; want 1", ans)
	}
}

func TestAPID(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000011111111111,
		0,
		0}

	ans := sph.APID()
	if ans != 0b11111111111 {
		t.Errorf("APID() = %b; want 11111111111", ans)
	}
}

func TestGroupingFlags(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b1100000000000000,
		0}

	ans := sph.GroupingFlags()
	if ans != 0b11 {
		t.Errorf("GroupingFlags() = %02b; want 11", ans)
	}
}

func TestSequenceCount(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b0011111111111111,
		0}

	ans := sph.SequenceCount()
	if ans != 0b11111111111111 {
		t.Errorf("SequenceCount() = %02b; want 11111111111111", ans)
	}
}
