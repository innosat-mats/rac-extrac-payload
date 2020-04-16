package innosat

import "testing"

func TestTMDataFieldHeader_PUSVersion(t *testing.T) {
	tests := []struct {
		name string
		h    *TMDataFieldHeader
		want uint8
	}{
		{
			"bitpattern",
			&TMDataFieldHeader{0b01110000, 0, 0, 0, 0},
			0b111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.PUSVersion(); got != tt.want {
				t.Errorf("TMDataFieldHeader.PUSVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
