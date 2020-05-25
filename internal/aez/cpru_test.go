package aez

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCPRU_Report(t *testing.T) {
	type fields struct {
		VGATE0 gate
		VSUBS0 subs
		VRD0   rd
		VOD0   od
		VGATE1 gate
		VSUBS1 subs
		VRD1   rd
		VOD1   od
		VGATE2 gate
		VSUBS2 subs
		VRD2   rd
		VOD2   od
		VGATE3 gate
		VSUBS3 subs
		VRD3   rd
		VOD3   od
	}
	gate10 := gate(10)
	subs10 := subs(10)
	rd10 := rd(10)
	od10 := od(10)
	tests := []struct {
		name   string
		fields fields
		want   CPRUReport
	}{
		{"Transforms VGATE0", fields{VGATE0: 10}, CPRUReport{VGATE0: gate10.voltage()}},
		{"Transforms VSUBS0", fields{VSUBS0: 10}, CPRUReport{VSUBS0: subs10.voltage()}},
		{"Transforms VRD0", fields{VRD0: 10}, CPRUReport{VRD0: rd10.voltage()}},
		{"Transforms VOD0", fields{VOD0: 10}, CPRUReport{VOD0: od10.voltage()}},
		{"Transforms VGATE1", fields{VGATE1: 10}, CPRUReport{VGATE1: gate10.voltage()}},
		{"Transforms VSUBS1", fields{VSUBS1: 10}, CPRUReport{VSUBS1: subs10.voltage()}},
		{"Transforms VRD1", fields{VRD1: 10}, CPRUReport{VRD1: rd10.voltage()}},
		{"Transforms VOD1", fields{VOD1: 10}, CPRUReport{VOD1: od10.voltage()}},
		{"Transforms VGATE2", fields{VGATE2: 10}, CPRUReport{VGATE2: gate10.voltage()}},
		{"Transforms VSUBS2", fields{VSUBS2: 10}, CPRUReport{VSUBS2: subs10.voltage()}},
		{"Transforms VRD2", fields{VRD2: 10}, CPRUReport{VRD2: rd10.voltage()}},
		{"Transforms VOD2", fields{VOD2: 10}, CPRUReport{VOD2: od10.voltage()}},
		{"Transforms VGATE3", fields{VGATE3: 10}, CPRUReport{VGATE3: gate10.voltage()}},
		{"Transforms VSUBS3", fields{VSUBS3: 10}, CPRUReport{VSUBS3: subs10.voltage()}},
		{"Transforms VRD3", fields{VRD3: 10}, CPRUReport{VRD3: rd10.voltage()}},
		{"Transforms VOD3", fields{VOD3: 10}, CPRUReport{VOD3: od10.voltage()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpru := &CPRU{
				VGATE0: tt.fields.VGATE0,
				VSUBS0: tt.fields.VSUBS0,
				VRD0:   tt.fields.VRD0,
				VOD0:   tt.fields.VOD0,
				VGATE1: tt.fields.VGATE1,
				VSUBS1: tt.fields.VSUBS1,
				VRD1:   tt.fields.VRD1,
				VOD1:   tt.fields.VOD1,
				VGATE2: tt.fields.VGATE2,
				VSUBS2: tt.fields.VSUBS2,
				VRD2:   tt.fields.VRD2,
				VOD2:   tt.fields.VOD2,
				VGATE3: tt.fields.VGATE3,
				VSUBS3: tt.fields.VSUBS3,
				VRD3:   tt.fields.VRD3,
				VOD3:   tt.fields.VOD3,
			}
			if got := cpru.Report(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CPRU.Report() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPRU_CSVHeaders(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			"Produces headers",
			[]string{
				"VGATE0", "VSUBS0", "VRD0", "VOD0", "Overvoltage0", "Power0",
				"VGATE1", "VSUBS1", "VRD1", "VOD1", "Overvoltage1", "Power1",
				"VGATE2", "VSUBS2", "VRD2", "VOD2", "Overvoltage2", "Power2",
				"VGATE3", "VSUBS3", "VRD3", "VOD3", "Overvoltage3", "Power3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpru := CPRU{}
			if got := cpru.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CPRU.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPRU_CSVRow(t *testing.T) {
	type fields struct {
		STAT   cpruStat
		VGATE0 gate
		VSUBS0 subs
		VRD0   rd
		VOD0   od
		VGATE1 gate
		VSUBS1 subs
		VRD1   rd
		VOD1   od
		VGATE2 gate
		VSUBS2 subs
		VRD2   rd
		VOD2   od
		VGATE3 gate
		VSUBS3 subs
		VRD3   rd
		VOD3   od
	}
	gate2 := gate(2)
	subs3 := subs(3)
	rd4 := rd(4)
	od5 := od(5)
	cpruStat1 := cpruStat(1)
	gate6 := gate(6)
	subs7 := subs(7)
	rd8 := rd(8)
	od9 := od(9)
	gate10 := gate(10)
	subs11 := subs(11)
	rd12 := rd(12)
	od13 := od(13)
	gate14 := gate(14)
	subs15 := subs(15)
	rd16 := rd(16)
	od17 := od(17)
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates data row",
			fields{
				STAT:   1,
				VGATE0: 2, VSUBS0: 3, VRD0: 4, VOD0: 5,
				VGATE1: 6, VSUBS1: 7, VRD1: 8, VOD1: 9,
				VGATE2: 10, VSUBS2: 11, VRD2: 12, VOD2: 13,
				VGATE3: 14, VSUBS3: 15, VRD3: 16, VOD3: 17,
			},
			[]string{
				fmt.Sprintf("%v", gate2.voltage()),
				fmt.Sprintf("%v", subs3.voltage()),
				fmt.Sprintf("%v", rd4.voltage()),
				fmt.Sprintf("%v", od5.voltage()),
				fmt.Sprintf("%v", cpruStat1.overvoltageFault(0)),
				fmt.Sprintf("%v", cpruStat1.powerEnabled(0)),
				fmt.Sprintf("%v", gate6.voltage()),
				fmt.Sprintf("%v", subs7.voltage()),
				fmt.Sprintf("%v", rd8.voltage()),
				fmt.Sprintf("%v", od9.voltage()),
				fmt.Sprintf("%v", cpruStat1.overvoltageFault(1)),
				fmt.Sprintf("%v", cpruStat1.powerEnabled(1)),
				fmt.Sprintf("%v", gate10.voltage()),
				fmt.Sprintf("%v", subs11.voltage()),
				fmt.Sprintf("%v", rd12.voltage()),
				fmt.Sprintf("%v", od13.voltage()),
				fmt.Sprintf("%v", cpruStat1.overvoltageFault(2)),
				fmt.Sprintf("%v", cpruStat1.powerEnabled(2)),
				fmt.Sprintf("%v", gate14.voltage()),
				fmt.Sprintf("%v", subs15.voltage()),
				fmt.Sprintf("%v", rd16.voltage()),
				fmt.Sprintf("%v", od17.voltage()),
				fmt.Sprintf("%v", cpruStat1.overvoltageFault(3)),
				fmt.Sprintf("%v", cpruStat1.powerEnabled(3)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpru := CPRU{
				STAT:   tt.fields.STAT,
				VGATE0: tt.fields.VGATE0,
				VSUBS0: tt.fields.VSUBS0,
				VRD0:   tt.fields.VRD0,
				VOD0:   tt.fields.VOD0,
				VGATE1: tt.fields.VGATE1,
				VSUBS1: tt.fields.VSUBS1,
				VRD1:   tt.fields.VRD1,
				VOD1:   tt.fields.VOD1,
				VGATE2: tt.fields.VGATE2,
				VSUBS2: tt.fields.VSUBS2,
				VRD2:   tt.fields.VRD2,
				VOD2:   tt.fields.VOD2,
				VGATE3: tt.fields.VGATE3,
				VSUBS3: tt.fields.VSUBS3,
				VRD3:   tt.fields.VRD3,
				VOD3:   tt.fields.VOD3,
			}
			if got := cpru.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CPRU.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPRU_CSVSpecifications(t *testing.T) {
	type fields struct {
		STAT   cpruStat
		VGATE0 gate
		VSUBS0 subs
		VRD0   rd
		VOD0   od
		VGATE1 gate
		VSUBS1 subs
		VRD1   rd
		VOD1   od
		VGATE2 gate
		VSUBS2 subs
		VRD2   rd
		VOD2   od
		VGATE3 gate
		VSUBS3 subs
		VRD3   rd
		VOD3   od
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
			cpru := CPRU{
				STAT:   tt.fields.STAT,
				VGATE0: tt.fields.VGATE0,
				VSUBS0: tt.fields.VSUBS0,
				VRD0:   tt.fields.VRD0,
				VOD0:   tt.fields.VOD0,
				VGATE1: tt.fields.VGATE1,
				VSUBS1: tt.fields.VSUBS1,
				VRD1:   tt.fields.VRD1,
				VOD1:   tt.fields.VOD1,
				VGATE2: tt.fields.VGATE2,
				VSUBS2: tt.fields.VSUBS2,
				VRD2:   tt.fields.VRD2,
				VOD2:   tt.fields.VOD2,
				VGATE3: tt.fields.VGATE3,
				VSUBS3: tt.fields.VSUBS3,
				VRD3:   tt.fields.VRD3,
				VOD3:   tt.fields.VOD3,
			}
			if got := cpru.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CPRU.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cpruStat_overvoltageFault(t *testing.T) {
	type args struct {
		ccd uint8
	}
	tests := []struct {
		name string
		stat cpruStat
		args args
		want bool
	}{
		{"overvoltageFault true for CCD0", cpruStat(0x80), args{0}, true},
		{"overvoltageFault false for CCD0", cpruStat(0x7F), args{0}, false},
		{"overvoltageFault true for CCD1", cpruStat(0x40), args{1}, true},
		{"overvoltageFault false for CCD1", cpruStat(0xBF), args{1}, false},
		{"overvoltageFault true for CCD2", cpruStat(0x20), args{2}, true},
		{"overvoltageFault false for CCD2", cpruStat(0xDF), args{2}, false},
		{"overvoltageFault true for CCD3", cpruStat(0x10), args{3}, true},
		{"overvoltageFault false for CCD3", cpruStat(0xEF), args{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stat.overvoltageFault(tt.args.ccd); got != tt.want {
				t.Errorf("cpruStat.overvoltageFault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cpruStat_powerEnabled(t *testing.T) {
	type args struct {
		ccd uint8
	}
	tests := []struct {
		name string
		stat cpruStat
		args args
		want bool
	}{
		{"powerEnabled true for CCD0", cpruStat(0x08), args{0}, true},
		{"powerEnabled false for CCD0", cpruStat(0xF7), args{0}, false},
		{"powerEnabled true for CCD1", cpruStat(0x04), args{1}, true},
		{"powerEnabled false for CCD1", cpruStat(0xFB), args{1}, false},
		{"powerEnabled true for CCD2", cpruStat(0x02), args{2}, true},
		{"powerEnabled false for CCD2", cpruStat(0xFD), args{2}, false},
		{"powerEnabled true for CCD3", cpruStat(0x01), args{3}, true},
		{"powerEnabled false for CCD3", cpruStat(0xFE), args{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stat.powerEnabled(tt.args.ccd); got != tt.want {
				t.Errorf("cpruStat.powerEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
