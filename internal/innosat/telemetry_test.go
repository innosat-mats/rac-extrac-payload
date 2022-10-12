package innosat

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

func TestTMHeader_PUSVersion(t *testing.T) {
	tests := []struct {
		name string
		h    *TMHeader
		want uint8
	}{
		{
			"bitpattern",
			&TMHeader{0b01110000, 0, 0, 0, 0},
			0b111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.PUS.Version(); got != tt.want {
				t.Errorf("TMHeader.PUSVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_Time(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
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
			h := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := h.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMHeader.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_Nanoseconds(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
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
			h := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := h.Nanoseconds(); got != tt.want {
				t.Errorf("TMHeader.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_IsHousekeeping(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
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
			tmdfh := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.IsHousekeeping(); got != tt.want {
				t.Errorf("TMHeader.IsHousekeeping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_IsTransparentData(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
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
			tmdfh := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.IsTransparentData(); got != tt.want {
				t.Errorf("TMHeader.IsTransparentData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_CSVHeaders(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Generates headers", fields{}, []string{"TMHeaderTime", "TMHeaderNanoseconds"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMHeader.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_CSVRow(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a data row",
			fields{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
			[]string{"1980-01-06T00:00:24.75Z", "42750000000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMHeader.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_MarshalJSON(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"Marshals into expected json",
			fields{CUCTimeSeconds: 4},
			[]byte("{\"tmHeaderTime\":\"1980-01-05T23:59:46Z\",\"tmHeaderNanoseconds\":4000000000}"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			got, err := tmdfh.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TMHeader.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TMHeader.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTMHeader_IsTCVerification(t *testing.T) {
	type fields struct {
		PUS             pus
		ServiceType     SourcePackageServiceType
		ServiceSubType  SourcePackageServiceSubtype
		CUCTimeSeconds  uint32
		CUCTimeFraction uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Service type 1 and sub 1 is", fields{ServiceType: 1, ServiceSubType: 1}, true},
		{"Service type 1 and sub 2 is", fields{ServiceType: 1, ServiceSubType: 2}, true},
		{"Service type 1 and sub 7 is", fields{ServiceType: 1, ServiceSubType: 7}, true},
		{"Service type 1 and sub 8 is", fields{ServiceType: 1, ServiceSubType: 8}, true},
		{"Service type 1 and sub 42 is not", fields{ServiceType: 1, ServiceSubType: 42}, false},
		{"Service type 1 is not", fields{ServiceType: 1}, false},
		{"Service subtype 1 is not", fields{ServiceSubType: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdfh := &TMHeader{
				PUS:             tt.fields.PUS,
				ServiceType:     tt.fields.ServiceType,
				ServiceSubType:  tt.fields.ServiceSubType,
				CUCTimeSeconds:  tt.fields.CUCTimeSeconds,
				CUCTimeFraction: tt.fields.CUCTimeFraction,
			}
			if got := tmdfh.IsTCVerification(); got != tt.want {
				t.Errorf("TMHeader.IsTCVerification() = %v, want %v", got, tt.want)
			}
		})
	}
}
