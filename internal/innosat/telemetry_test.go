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
		ServiceType     SourcePackageServiceType
		ServiceSubType  uint8
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	type args struct {
		epoch time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			"Returns Epoch/TAI",
			fields{},
			args{ccsds.TAI},
			ccsds.TAI,
		},
		{
			"Returns expected time",
			fields{CUCTimeSeconds: 10, CUCTimeFraction: 0b1000000000000000},
			args{ccsds.TAI},
			ccsds.TAI.Add(time.Second * 10).Add(time.Millisecond * 500),
		},
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
			if got := h.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMDataFieldHeader.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMDataFieldHeader_Nanoseconds(t *testing.T) {
	type fields struct {
		PUS             uint8
		ServiceType     SourcePackageServiceType
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
		{"Returns nanoseconds", fields{0, 0, 0, 10, 0b1000000000000000}, 10500000000},
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

func TestTMDataFieldHeader_IsHousekeeping(t *testing.T) {
	type fields struct {
		PUS             uint8
		ServiceType     SourcePackageServiceType
		ServiceSubType  uint8
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"False in general", fields{}, false},
		{"False if only correct ServiceType", fields{ServiceType: HousekeepingDiagnosticDataReporting}, false},
		{"False if only correct ServiceSubType", fields{ServiceSubType: 25}, false},
		{
			"True if correct ServiceType and ServiceSubType",
			fields{ServiceType: HousekeepingDiagnosticDataReporting, ServiceSubType: 25},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := &TMDataFieldHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.IsHousekeeping(); got != tt.want {
				t.Errorf("TMDataFieldHeader.IsHousekeeping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMDataFieldHeader_IsTransparentData(t *testing.T) {
	type fields struct {
		PUS             uint8
		ServiceType     SourcePackageServiceType
		ServiceSubType  uint8
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Service type 128 and sub 25 is", fields{ServiceType: 128, ServiceSubType: 25}, true},
		{"ServiceSubType 25 is not", fields{ServiceSubType: 25}, false},
		{"Service type 128 is not", fields{ServiceType: 128}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := &TMDataFieldHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.IsTransparentData(); got != tt.want {
				t.Errorf("TMDataFieldHeader.IsTransparentData() = %v, want %v", got, tt.want)
			}
		})
	}
}
