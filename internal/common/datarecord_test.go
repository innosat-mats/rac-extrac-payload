package common

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func TestDataRecord_CSVSpecifications(t *testing.T) {
	type fields struct {
		Origin       OriginDescription
		RamsesHeader ramses.Ramses
		RamsesSecure ramses.Secure
		SourceHeader innosat.SourcePacketHeader
		TMHeader     innosat.TMDataFieldHeader
		SID          aez.SID
		Data         Exportable
		Error        error
		Buffer       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Works on empty package",
			fields{},
			[]string{"RAMSES", ramses.Specification, "INNOSAT", innosat.Specification},
		},
		{
			"Returns specs",
			fields{
				RamsesHeader: ramses.Ramses{Type: 4},
				SourceHeader: innosat.SourcePacketHeader{PacketLength: 42},
				Data:         aez.HTR{},
			},
			[]string{
				"RAMSES", ramses.Specification,
				"INNOSAT", innosat.Specification,
				"AEZ", aez.Specification,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := DataRecord{
				Origin:       tt.fields.Origin,
				RamsesHeader: tt.fields.RamsesHeader,
				RamsesSecure: tt.fields.RamsesSecure,
				SourceHeader: tt.fields.SourceHeader,
				TMHeader:     tt.fields.TMHeader,
				SID:          tt.fields.SID,
				Data:         tt.fields.Data,
				Error:        tt.fields.Error,
				Buffer:       tt.fields.Buffer,
			}
			got := record.CSVSpecifications()

			if (len(got) > 0 || len(tt.want) > 0) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRecord.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRecord_CSVHeaders(t *testing.T) {
	type fields struct {
		Origin       OriginDescription
		RamsesHeader ramses.Ramses
		RamsesSecure ramses.Secure
		SourceHeader innosat.SourcePacketHeader
		TMHeader     innosat.TMDataFieldHeader
		SID          aez.SID
		Data         Exportable
		Error        error
		Buffer       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Handles missing Data",
			fields{},
			[]string{
				"File",
				"ProcessingDate",
				"RamsesTime",
				"SPSequenceCount",
				"TMHeaderTime",
				"TMHeaderNanoseconds",
				"SID",
				"Error",
			},
		},
		{
			"Returns expected headers",
			fields{Data: aez.STAT{}},
			[]string{
				"File",
				"ProcessingDate",
				"RamsesTime",
				"SPSequenceCount",
				"TMHeaderTime",
				"TMHeaderNanoseconds",
				"SID",
				"STATTIME",
				"STATNANO",
				"SPID",
				"SPREV",
				"FPID",
				"FPREV",
				"SVNA",
				"SVNB",
				"SVNC",
				"MODE",
				"EDACE",
				"EDACCE",
				"EDACN",
				"SPWEOP",
				"SPWEEP",
				"ANOMALY",
				"Error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := DataRecord{
				Origin:       tt.fields.Origin,
				RamsesHeader: tt.fields.RamsesHeader,
				RamsesSecure: tt.fields.RamsesSecure,
				SourceHeader: tt.fields.SourceHeader,
				TMHeader:     tt.fields.TMHeader,
				SID:          tt.fields.SID,
				Data:         tt.fields.Data,
				Error:        tt.fields.Error,
				Buffer:       tt.fields.Buffer,
			}
			if got := record.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRecord.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

var procDate = time.Now()

func TestDataRecord_CSVRow(t *testing.T) {
	type fields struct {
		Origin       OriginDescription
		RamsesHeader ramses.Ramses
		RamsesSecure ramses.Secure
		SourceHeader innosat.SourcePacketHeader
		TMHeader     innosat.TMDataFieldHeader
		SID          aez.SID
		Data         Exportable
		Error        error
		Buffer       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Handles missing Data",
			fields{
				Origin:       OriginDescription{Name: "Sputnik", ProcessingDate: procDate},
				RamsesHeader: ramses.Ramses{Date: 24, Time: 42000},
				SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: innosat.PacketSequenceControl(0xc003)},
				TMHeader:     innosat.TMDataFieldHeader{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
				SID:          aez.SIDSTAT,
				Error:        errors.New("Test"),
			},
			[]string{
				"Sputnik",
				procDate.Format(time.RFC3339),
				"2000-01-25T00:00:42Z",
				"3",
				"1980-01-06T00:00:42.75Z",
				"42750000000",
				"STAT",
				"Test",
			},
		},
		{
			"Handles missing Error",
			fields{
				Origin:       OriginDescription{Name: "Sputnik", ProcessingDate: procDate},
				RamsesHeader: ramses.Ramses{Date: 24, Time: 42000},
				SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: innosat.PacketSequenceControl(0xc003)},
				TMHeader:     innosat.TMDataFieldHeader{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
				SID:          aez.SIDSTAT,
				Data: aez.STAT{
					SPID:    1,
					SPREV:   2,
					FPID:    3,
					FPREV:   4,
					SVNA:    5,
					SVNB:    6,
					SVNC:    7,
					TS:      8,
					TSS:     9,
					MODE:    10,
					EDACE:   11,
					EDACCE:  12,
					EDACN:   13,
					SPWEOP:  14,
					SPWEEP:  15,
					ANOMALY: 16,
				},
			},
			[]string{
				"Sputnik",
				procDate.Format(time.RFC3339),
				"2000-01-25T00:00:42Z",
				"3",
				"1980-01-06T00:00:42.75Z",
				"42750000000",
				"STAT",
				"1980-01-06T00:00:08.000137329Z",
				"8000137329",
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"10",
				"11",
				"12",
				"13",
				"14",
				"15",
				"16",
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := DataRecord{
				Origin:       tt.fields.Origin,
				RamsesHeader: tt.fields.RamsesHeader,
				RamsesSecure: tt.fields.RamsesSecure,
				SourceHeader: tt.fields.SourceHeader,
				TMHeader:     tt.fields.TMHeader,
				SID:          tt.fields.SID,
				Data:         tt.fields.Data,
				Error:        tt.fields.Error,
				Buffer:       tt.fields.Buffer,
			}
			if got := record.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRecord.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRecord_AEZData(t *testing.T) {
	type fields struct {
		Origin       OriginDescription
		RamsesHeader ramses.Ramses
		RamsesSecure ramses.Secure
		SourceHeader innosat.SourcePacketHeader
		TMHeader     innosat.TMDataFieldHeader
		SID          aez.SID
		Data         Exportable
		Error        error
		Buffer       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"Handles Data not set", fields{}, nil},
		{"Returns Data", fields{Data: aez.CPRU{}}, aez.CPRU{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := DataRecord{
				Origin:       tt.fields.Origin,
				RamsesHeader: tt.fields.RamsesHeader,
				RamsesSecure: tt.fields.RamsesSecure,
				SourceHeader: tt.fields.SourceHeader,
				TMHeader:     tt.fields.TMHeader,
				SID:          tt.fields.SID,
				Data:         tt.fields.Data,
				Error:        tt.fields.Error,
				Buffer:       tt.fields.Buffer,
			}
			if got := record.AEZData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRecord.AEZData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRecord_OriginName(t *testing.T) {
	type fields struct {
		Origin       OriginDescription
		RamsesHeader ramses.Ramses
		RamsesSecure ramses.Secure
		SourceHeader innosat.SourcePacketHeader
		TMHeader     innosat.TMDataFieldHeader
		SID          aez.SID
		Data         Exportable
		Error        error
		Buffer       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Returns the Origin.Name", fields{Origin: OriginDescription{Name: "heeelloo"}}, "heeelloo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := DataRecord{
				Origin:       tt.fields.Origin,
				RamsesHeader: tt.fields.RamsesHeader,
				RamsesSecure: tt.fields.RamsesSecure,
				SourceHeader: tt.fields.SourceHeader,
				TMHeader:     tt.fields.TMHeader,
				SID:          tt.fields.SID,
				Data:         tt.fields.Data,
				Error:        tt.fields.Error,
				Buffer:       tt.fields.Buffer,
			}
			if got := record.OriginName(); got != tt.want {
				t.Errorf("DataRecord.OriginName() = %v, want %v", got, tt.want)
			}
		})
	}
}
