package timeseries

import (
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func TestOutStream_String(t *testing.T) {
	tests := []struct {
		name   string
		stream OutStream
		want   string
	}{
		{"HTR", HTR, "HTR"},
		{"PWR", PWR, "PWR"},
		{"CPRU", CPRU, "CPRU"},
		{"STAT", STAT, "STAT"},
		{"PM", PM, "PM"},
		{"CCD", CCD, "CCD"},
		{"TCV", TCV, "TCV"},
		{"default", Unknown, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stream.String(); got != tt.want {
				t.Errorf("OutStream.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutStreamFromDataRecord(t *testing.T) {
	type args struct {
		pkg *common.DataRecord
	}
	tests := []struct {
		name string
		args args
		want OutStream
	}{
		{"CCD", args{&common.DataRecord{Data: &aez.CCDImage{}}}, CCD},
		{"PM", args{&common.DataRecord{Data: &aez.PMData{}}}, PM},
		{"HTR", args{&common.DataRecord{Data: &aez.HTR{}}}, HTR},
		{"PWR", args{&common.DataRecord{Data: &aez.PWR{}}}, PWR},
		{"CPRU", args{&common.DataRecord{Data: &aez.CPRU{}}}, CPRU},
		{"STAT", args{&common.DataRecord{Data: &aez.STAT{}}}, STAT},
		{"TCV, accept success", args{&common.DataRecord{Data: &aez.TCAcceptSuccessData{}}}, TCV},
		{"TCV, accept fail", args{&common.DataRecord{Data: &aez.TCAcceptFailureData{}}}, TCV},
		{"TCV, exec success", args{&common.DataRecord{Data: &aez.TCExecSuccessData{}}}, TCV},
		{"TCV, exec fail", args{&common.DataRecord{Data: &aez.TCExecFailureData{}}}, TCV},
		{"Unknown", args{&common.DataRecord{}}, Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OutStreamFromDataRecord(tt.args.pkg); got != tt.want {
				t.Errorf("OutStreamFromDataRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
