package common

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func TestDataRecord_CSVSpecifications(t *testing.T) {
	type fields struct {
		Origin         OriginDescription
		RamsesHeader   ramses.Ramses
		RamsesTMHeader ramses.TMHeader
		SourceHeader   innosat.SourcePacketHeader
		TMHeader       innosat.TMHeader
		SID            aez.SID
		RID            aez.RID
		Data           Exporter
		Error          error
		Buffer         []byte
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
				Origin:         tt.fields.Origin,
				RamsesHeader:   tt.fields.RamsesHeader,
				RamsesTMHeader: tt.fields.RamsesTMHeader,
				SourceHeader:   tt.fields.SourceHeader,
				TMHeader:       tt.fields.TMHeader,
				SID:            tt.fields.SID,
				RID:            tt.fields.RID,
				Data:           tt.fields.Data,
				Error:          tt.fields.Error,
				Buffer:         tt.fields.Buffer,
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
		Origin         OriginDescription
		RamsesHeader   ramses.Ramses
		RamsesTMHeader ramses.TMHeader
		SourceHeader   innosat.SourcePacketHeader
		TMHeader       innosat.TMHeader
		SID            aez.SID
		RID            aez.RID
		Data           Exporter
		Error          error
		Buffer         []byte
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
				"QualityIndicator",
				"LossFlag",
				"VCFrameCounter",
				"SPSequenceCount",
				"TMHeaderTime",
				"TMHeaderNanoseconds",
				"SID",
				"RID",
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
				"QualityIndicator",
				"LossFlag",
				"VCFrameCounter",
				"SPSequenceCount",
				"TMHeaderTime",
				"TMHeaderNanoseconds",
				"SID",
				"RID",
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
				Origin:         tt.fields.Origin,
				RamsesHeader:   tt.fields.RamsesHeader,
				RamsesTMHeader: tt.fields.RamsesTMHeader,
				SourceHeader:   tt.fields.SourceHeader,
				TMHeader:       tt.fields.TMHeader,
				SID:            tt.fields.SID,
				RID:            tt.fields.RID,
				Data:           tt.fields.Data,
				Error:          tt.fields.Error,
				Buffer:         tt.fields.Buffer,
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
		Origin         OriginDescription
		RamsesHeader   ramses.Ramses
		RamsesTMHeader ramses.TMHeader
		SourceHeader   innosat.SourcePacketHeader
		TMHeader       innosat.TMHeader
		SID            aez.SID
		RID            aez.RID
		Data           Exporter
		Error          error
		Buffer         []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Handles missing Data",
			fields{
				Origin:         OriginDescription{Name: "Sputnik", ProcessingDate: procDate},
				RamsesHeader:   ramses.Ramses{Date: 24, Time: 42000},
				RamsesTMHeader: ramses.TMHeader{LossFlag: 1, VCFrameCounter: 42},
				SourceHeader:   innosat.SourcePacketHeader{PacketSequenceControl: innosat.PacketSequenceControl(0xc003)},
				TMHeader:       innosat.TMHeader{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
				SID:            aez.SIDSTAT,
				RID:            aez.CCD1,
				Error:          errors.New("Test"),
			},
			[]string{
				"Sputnik",
				procDate.Format(time.RFC3339),
				"2000-01-25T00:00:42Z",
				"0",
				"1",
				"42",
				"3",
				"1980-01-06T00:00:42.75Z",
				"42750000000",
				"STAT",
				"CCD1",
				"Test",
			},
		},
		{
			"Handles missing Error",
			fields{
				Origin:         OriginDescription{Name: "Sputnik", ProcessingDate: procDate},
				RamsesHeader:   ramses.Ramses{Date: 24, Time: 42000},
				RamsesTMHeader: ramses.TMHeader{LossFlag: 1, VCFrameCounter: 42},
				SourceHeader:   innosat.SourcePacketHeader{PacketSequenceControl: innosat.PacketSequenceControl(0xc003)},
				TMHeader:       innosat.TMHeader{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
				SID:            aez.SIDSTAT,
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
				"0",
				"1",
				"42",
				"3",
				"1980-01-06T00:00:42.75Z",
				"42750000000",
				"STAT",
				"",
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
				Origin:         tt.fields.Origin,
				RamsesHeader:   tt.fields.RamsesHeader,
				RamsesTMHeader: tt.fields.RamsesTMHeader,
				SourceHeader:   tt.fields.SourceHeader,
				TMHeader:       tt.fields.TMHeader,
				SID:            tt.fields.SID,
				RID:            tt.fields.RID,
				Data:           tt.fields.Data,
				Error:          tt.fields.Error,
				Buffer:         tt.fields.Buffer,
			}
			if got := record.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRecord.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRecord_MarshalJSON(t *testing.T) {
	type fields struct {
		Data  Exporter
		Error error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"No Data, No Error", fields{}},
		{"Error", fields{Error: io.EOF}},
		{"Image Data", fields{Data: &aez.CCDImage{PackData: &aez.CCDImagePackData{}}}},
		{"Non-image Data", fields{Data: aez.STAT{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &DataRecord{
				Data:  tt.fields.Data,
				Error: tt.fields.Error,
			}
			got, err := record.MarshalJSON()
			if err != nil {
				t.Errorf("DataRecord.MarshalJSON() error = %v", err)
				return
			}
			var js map[string]interface{}
			if json.Unmarshal(got, &js) != nil {
				t.Errorf("DataRecord.MarshalJSON() = %v, not a valid json", string(got))
			}
		})
	}
}
