package ramses

import (
	"reflect"
	"testing"
)

func TestTMHeader_CSVHeader(t *testing.T) {
	header := TMHeader{}
	want := []string{"QualityIndicator", "LossFlag", "VCFrameCounter"}
	if got := header.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("TMHeader.CSVHeader() = %v, wnat %v", got, want)
	}
}

func TestTMHeader_CSVRow(t *testing.T) {
	header := TMHeader{
		QualityIndicator: CompletePacket,
		LossFlag:         Discontinuities,
		VCFrameCounter:   42,
	}

	want := []string{"0", "1", "42"}

	if got := header.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("TMHeader.CSVROW() = %v, want %v", got, want)
	}
}
