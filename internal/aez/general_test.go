package aez

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

func TestSTAT_Time(t *testing.T) {
	type fields struct {
		SPID   uint16
		SPREV  uint8
		FPID   uint16
		FPREV  uint8
		TS     uint32
		TSS    uint16
		MODE   uint8
		EDACE  uint32
		EDACCE uint32
		EDACN  uint32
		SPWEOP uint32
		SPWEEP uint32
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
				SPID:   tt.fields.SPID,
				SPREV:  tt.fields.SPREV,
				FPID:   tt.fields.FPID,
				FPREV:  tt.fields.FPREV,
				TS:     tt.fields.TS,
				TSS:    tt.fields.TSS,
				MODE:   tt.fields.MODE,
				EDACE:  tt.fields.EDACE,
				EDACCE: tt.fields.EDACCE,
				EDACN:  tt.fields.EDACN,
				SPWEOP: tt.fields.SPWEOP,
				SPWEEP: tt.fields.SPWEEP,
			}
			if got := stat.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STAT.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTAT_Nanoseconds(t *testing.T) {
	type fields struct {
		SPID   uint16
		SPREV  uint8
		FPID   uint16
		FPREV  uint8
		TS     uint32
		TSS    uint16
		MODE   uint8
		EDACE  uint32
		EDACCE uint32
		EDACN  uint32
		SPWEOP uint32
		SPWEEP uint32
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
				SPID:   tt.fields.SPID,
				SPREV:  tt.fields.SPREV,
				FPID:   tt.fields.FPID,
				FPREV:  tt.fields.FPREV,
				TS:     tt.fields.TS,
				TSS:    tt.fields.TSS,
				MODE:   tt.fields.MODE,
				EDACE:  tt.fields.EDACE,
				EDACCE: tt.fields.EDACCE,
				EDACN:  tt.fields.EDACN,
				SPWEOP: tt.fields.SPWEOP,
				SPWEEP: tt.fields.SPWEEP,
			}
			if got := stat.Nanoseconds(); got != tt.want {
				t.Errorf("STAT.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
