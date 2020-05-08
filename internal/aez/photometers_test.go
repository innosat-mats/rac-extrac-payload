package aez

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

func TestPMData_Time(t *testing.T) {
	type fields struct {
		EXPTS    uint32
		EXPTSS   uint16
		PM1A     uint32
		PM1ACNTR uint32
		PM1B     uint32
		PM1BCNTR uint32
		PM1S     uint32
		PM1SCNTR uint32
		PM2A     uint32
		PM2ACNTR uint32
		PM2B     uint32
		PM2BCNTR uint32
		PM2S     uint32
		PM2SCNTR uint32
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
			fields{EXPTS: 10, EXPTSS: 0xC000},
			args{ccsds.TAI},
			ccsds.TAI.Add(time.Second * 10).Add(time.Millisecond * 750),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PMData{
				EXPTS:    tt.fields.EXPTS,
				EXPTSS:   tt.fields.EXPTSS,
				PM1A:     tt.fields.PM1A,
				PM1ACNTR: tt.fields.PM1ACNTR,
				PM1B:     tt.fields.PM1B,
				PM1BCNTR: tt.fields.PM1BCNTR,
				PM1S:     tt.fields.PM1S,
				PM1SCNTR: tt.fields.PM1SCNTR,
				PM2A:     tt.fields.PM2A,
				PM2ACNTR: tt.fields.PM2ACNTR,
				PM2B:     tt.fields.PM2B,
				PM2BCNTR: tt.fields.PM2BCNTR,
				PM2S:     tt.fields.PM2S,
				PM2SCNTR: tt.fields.PM2SCNTR,
			}
			if got := pm.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PMData.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPMData_Nanoseconds(t *testing.T) {
	type fields struct {
		EXPTS    uint32
		EXPTSS   uint16
		PM1A     uint32
		PM1ACNTR uint32
		PM1B     uint32
		PM1BCNTR uint32
		PM1S     uint32
		PM1SCNTR uint32
		PM2A     uint32
		PM2ACNTR uint32
		PM2B     uint32
		PM2BCNTR uint32
		PM2S     uint32
		PM2SCNTR uint32
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
			fields{EXPTS: 10, EXPTSS: 0xC000},
			10750000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PMData{
				EXPTS:    tt.fields.EXPTS,
				EXPTSS:   tt.fields.EXPTSS,
				PM1A:     tt.fields.PM1A,
				PM1ACNTR: tt.fields.PM1ACNTR,
				PM1B:     tt.fields.PM1B,
				PM1BCNTR: tt.fields.PM1BCNTR,
				PM1S:     tt.fields.PM1S,
				PM1SCNTR: tt.fields.PM1SCNTR,
				PM2A:     tt.fields.PM2A,
				PM2ACNTR: tt.fields.PM2ACNTR,
				PM2B:     tt.fields.PM2B,
				PM2BCNTR: tt.fields.PM2BCNTR,
				PM2S:     tt.fields.PM2S,
				PM2SCNTR: tt.fields.PM2SCNTR,
			}
			if got := pm.Nanoseconds(); got != tt.want {
				t.Errorf("PMData.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPMData_CSVHeaders(t *testing.T) {
	type fields struct {
		EXPTS    uint32
		EXPTSS   uint16
		PM1A     uint32
		PM1ACNTR uint32
		PM1B     uint32
		PM1BCNTR uint32
		PM1S     uint32
		PM1SCNTR uint32
		PM2A     uint32
		PM2ACNTR uint32
		PM2B     uint32
		PM2BCNTR uint32
		PM2S     uint32
		PM2SCNTR uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates headers",
			fields{},
			[]string{
				"PMTIME",
				"PMNANO",
				"PM1A",
				"PM1ACNTR",
				"PM1B",
				"PM1BCNTR",
				"PM1S",
				"PM1SCNTR",
				"PM2A",
				"PM2ACNTR",
				"PM2B",
				"PM2BCNTR",
				"PM2S",
				"PM2SCNTR",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := PMData{
				EXPTS:    tt.fields.EXPTS,
				EXPTSS:   tt.fields.EXPTSS,
				PM1A:     tt.fields.PM1A,
				PM1ACNTR: tt.fields.PM1ACNTR,
				PM1B:     tt.fields.PM1B,
				PM1BCNTR: tt.fields.PM1BCNTR,
				PM1S:     tt.fields.PM1S,
				PM1SCNTR: tt.fields.PM1SCNTR,
				PM2A:     tt.fields.PM2A,
				PM2ACNTR: tt.fields.PM2ACNTR,
				PM2B:     tt.fields.PM2B,
				PM2BCNTR: tt.fields.PM2BCNTR,
				PM2S:     tt.fields.PM2S,
				PM2SCNTR: tt.fields.PM2SCNTR,
			}
			if got := pm.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PMData.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPMData_CSVRow(t *testing.T) {
	type fields struct {
		EXPTS    uint32
		EXPTSS   uint16
		PM1A     uint32
		PM1ACNTR uint32
		PM1B     uint32
		PM1BCNTR uint32
		PM1S     uint32
		PM1SCNTR uint32
		PM2A     uint32
		PM2ACNTR uint32
		PM2B     uint32
		PM2BCNTR uint32
		PM2S     uint32
		PM2SCNTR uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{
				EXPTS:    1,
				EXPTSS:   2,
				PM1A:     3,
				PM1ACNTR: 4,
				PM1B:     5,
				PM1BCNTR: 6,
				PM1S:     7,
				PM1SCNTR: 8,
				PM2A:     9,
				PM2ACNTR: 10,
				PM2B:     11,
				PM2BCNTR: 12,
				PM2S:     13,
				PM2SCNTR: 14,
			},
			[]string{
				"1980-01-06T00:00:01.000030518Z",
				"1000030518",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
				"10",
				"11",
				"12",
				"13",
				"14",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := PMData{
				EXPTS:    tt.fields.EXPTS,
				EXPTSS:   tt.fields.EXPTSS,
				PM1A:     tt.fields.PM1A,
				PM1ACNTR: tt.fields.PM1ACNTR,
				PM1B:     tt.fields.PM1B,
				PM1BCNTR: tt.fields.PM1BCNTR,
				PM1S:     tt.fields.PM1S,
				PM1SCNTR: tt.fields.PM1SCNTR,
				PM2A:     tt.fields.PM2A,
				PM2ACNTR: tt.fields.PM2ACNTR,
				PM2B:     tt.fields.PM2B,
				PM2BCNTR: tt.fields.PM2BCNTR,
				PM2S:     tt.fields.PM2S,
				PM2SCNTR: tt.fields.PM2SCNTR,
			}
			if got := pm.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PMData.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPMData_CSVSpecifications(t *testing.T) {
	type fields struct {
		EXPTS    uint32
		EXPTSS   uint16
		PM1A     uint32
		PM1ACNTR uint32
		PM1B     uint32
		PM1BCNTR uint32
		PM1S     uint32
		PM1SCNTR uint32
		PM2A     uint32
		PM2ACNTR uint32
		PM2B     uint32
		PM2BCNTR uint32
		PM2S     uint32
		PM2SCNTR uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Genereates spec", fields{}, []string{"AEZ", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := PMData{
				EXPTS:    tt.fields.EXPTS,
				EXPTSS:   tt.fields.EXPTSS,
				PM1A:     tt.fields.PM1A,
				PM1ACNTR: tt.fields.PM1ACNTR,
				PM1B:     tt.fields.PM1B,
				PM1BCNTR: tt.fields.PM1BCNTR,
				PM1S:     tt.fields.PM1S,
				PM1SCNTR: tt.fields.PM1SCNTR,
				PM2A:     tt.fields.PM2A,
				PM2ACNTR: tt.fields.PM2ACNTR,
				PM2B:     tt.fields.PM2B,
				PM2BCNTR: tt.fields.PM2BCNTR,
				PM2S:     tt.fields.PM2S,
				PM2SCNTR: tt.fields.PM2SCNTR,
			}
			if got := pm.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PMData.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}
