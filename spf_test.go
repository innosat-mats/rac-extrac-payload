package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	sph := SourcePacketHeader{
		0b1110000000000000,
		0,
		0}

	ans := sph.getVersion()
	if sph.getVersion() != 0b111 {
		t.Errorf("getVersion() = %03b; want 111", ans)
	}
}

func TestType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0001000000000000,
		0,
		0}

	ans := sph.getType()
	if sph.getType() != 1 {
		t.Errorf("getType() = %b; want 1", ans)
	}
}
func TestHeaderType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000100000000000,
		0,
		0}

	ans := sph.getHeaderType()
	if sph.getHeaderType() != 1 {
		t.Errorf("getHeaderType() = %b; want 1", ans)
	}
}

func TestAPID(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000011111111111,
		0,
		0}

	ans := sph.getAPID()
	if sph.getAPID() != 0b11111111111 {
		t.Errorf("getAPID() = %b; want 11111111111", ans)
	}
}

func TestGroupingFlags(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b1100000000000000,
		0}

	ans := sph.getGroupingFlags()
	if sph.getGroupingFlags() != 0b11 {
		t.Errorf("getGroupingFlags() = %02b; want 11", ans)
	}
}

func TestSequenceCount(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b0011111111111111,
		0}

	ans := sph.getSequenceCount()
	if sph.getSequenceCount() != 0b11111111111111 {
		t.Errorf("getSequenceCount() = %02b; want 11111111111111", ans)
	}
}
