package innosat

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

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

func TestTMDataFieldHeader_Time(t *testing.T) {
	type fields struct {
		PUS             uint8
		ServiceType     uint8
		ServiceSubType  uint8
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"Returns Epoch/TAI", fields{0, 0, 0, 0, 0}, ccsds.TAI},
		{"Returns expected time", fields{0, 0, 0, 10, 2}, ccsds.TAI.Add(time.Second * 10).Add(time.Millisecond * 500)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TMDataFieldHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := h.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMDataFieldHeader.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMDataFieldHeader_Nanoseconds(t *testing.T) {
	type fields struct {
		PUS             uint8
		ServiceType     uint8
		ServiceSubType  uint8
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Returns 0", fields{0, 0, 0, 0, 0}, 0},
		{"Returns nanoseconds", fields{0, 0, 0, 10, 2}, 10500000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TMDataFieldHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := h.Nanoseconds(); got != tt.want {
				t.Errorf("TMDataFieldHeader.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
