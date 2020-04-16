package innosat

import "testing"

func TestTCDataFieldHeader_PUSVersion(t *testing.T) {
	tests := []struct {
		name string
		h    *TCDataFieldHeader
		want uint8
	}{
		{
			"bitpattern",
			&TCDataFieldHeader{0b01110000, 0, 0},
			0b111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.PUSVersion(); got != tt.want {
				t.Errorf("TCDataFieldHeader.PUSVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
