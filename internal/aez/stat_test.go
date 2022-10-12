package aez

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

func TestSTAT_Time(t *testing.T) {
	type fields struct {
		TS  uint32
		TSS uint16
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
		{"Defaults to epoch", fields{}, args{ccsds.TAI}, ccsds.TAI},
		{
			"Returns time after epoch",
			fields{TS: 10, TSS: 0b1100000000000000},
			args{ccsds.TAI},
			ccsds.TAI.Add(time.Second * 10).Add(time.Millisecond * 750),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := &STAT{
				TS:  tt.fields.TS,
				TSS: tt.fields.TSS,
			}
			if got := stat.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STAT.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTAT_Nanoseconds(t *testing.T) {
	type fields struct {
		TS  uint32
		TSS uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"Defaults to 0", fields{}, 0,
		},
		{
			"Returns nanoseconds",
			fields{TS: 10, TSS: 0b1100000000000000},
			10750000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := &STAT{
				TS:  tt.fields.TS,
				TSS: tt.fields.TSS,
			}
			if got := stat.Nanoseconds(); got != tt.want {
				t.Errorf("STAT.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTAT_CSVHeaders(t *testing.T) {
	want := []string{
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
	}
	stat := STAT{}
	if got := stat.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("STAT.CSVHeaders() = %v, want %v", got, want)
	}
}

func TestSTAT_CSVRow(t *testing.T) {
	type fields struct {
		SPID    uint16
		SPREV   uint8
		FPID    uint16
		FPREV   uint8
		SVNA    uint8
		SVNB    uint8
		SVNC    uint8
		TS      uint32
		TSS     uint16
		MODE    uint8
		EDACE   uint32
		EDACCE  uint32
		EDACN   uint32
		SPWEOP  uint32
		SPWEEP  uint32
		ANOMALY uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{
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
			[]string{
				"1980-01-05T23:59:50.000137329Z",
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := STAT{
				SPID:    tt.fields.SPID,
				SPREV:   tt.fields.SPREV,
				FPID:    tt.fields.FPID,
				FPREV:   tt.fields.FPREV,
				SVNA:    tt.fields.SVNA,
				SVNB:    tt.fields.SVNB,
				SVNC:    tt.fields.SVNC,
				TS:      tt.fields.TS,
				TSS:     tt.fields.TSS,
				MODE:    tt.fields.MODE,
				EDACE:   tt.fields.EDACE,
				EDACCE:  tt.fields.EDACCE,
				EDACN:   tt.fields.EDACN,
				SPWEOP:  tt.fields.SPWEOP,
				SPWEEP:  tt.fields.SPWEEP,
				ANOMALY: tt.fields.ANOMALY,
			}
			if got := stat.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STAT.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTAT_CSVSpecifications(t *testing.T) {
	stat := STAT{}
	want := []string{"AEZ", Specification}
	if got := stat.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("STAT.CSVSpecifications() = %v, want %v", got, want)
	}
}
