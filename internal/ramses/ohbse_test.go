package ramses

import (
	"reflect"
	"testing"
)

func TestOhbseCcsdsTMPacket_CSVHeaders(t *testing.T) {
	header := OhbseCcsdsTMPacket{}
	want := []string{"QualityIndicator", "LossFlag", "VCFrameCounter"}
	if got := header.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("OhbscCcsdsTMPacket.CSVHeader() = %v, wnat %v", got, want)
	}
}

func TestOhbseCcsdsTMPacket_CSVRow(t *testing.T) {
	header := OhbseCcsdsTMPacket{
		QualityIndicator: CompletePacket,
		LossFlag:         Discontinuities,
		VCFrameCounter:   42,
	}

	want := []string{"0", "1", "42"}

	if got := header.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("OhbscCcsdsTMPacket.CSVROW() = %v, want %v", got, want)
	}
}
