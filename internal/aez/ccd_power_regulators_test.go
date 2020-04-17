package aez

import (
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
	tests := []struct {
		name   string
		fields fields
		want   CPRUReport
	}{
		{"Transforms VGATE0", fields{VGATE0: 10}, CPRUReport{VGATE0: gate(10).voltage()}},
		{"Transforms VSUBS0", fields{VSUBS0: 10}, CPRUReport{VSUBS0: subs(10).voltage()}},
		{"Transforms VRD0", fields{VRD0: 10}, CPRUReport{VRD0: rd(10).voltage()}},
		{"Transforms VOD0", fields{VOD0: 10}, CPRUReport{VOD0: od(10).voltage()}},
		{"Transforms VGATE1", fields{VGATE1: 10}, CPRUReport{VGATE1: gate(10).voltage()}},
		{"Transforms VSUBS1", fields{VSUBS1: 10}, CPRUReport{VSUBS1: subs(10).voltage()}},
		{"Transforms VRD1", fields{VRD1: 10}, CPRUReport{VRD1: rd(10).voltage()}},
		{"Transforms VOD1", fields{VOD1: 10}, CPRUReport{VOD1: od(10).voltage()}},
		{"Transforms VGATE2", fields{VGATE2: 10}, CPRUReport{VGATE2: gate(10).voltage()}},
		{"Transforms VSUBS2", fields{VSUBS2: 10}, CPRUReport{VSUBS2: subs(10).voltage()}},
		{"Transforms VRD2", fields{VRD2: 10}, CPRUReport{VRD2: rd(10).voltage()}},
		{"Transforms VOD2", fields{VOD2: 10}, CPRUReport{VOD2: od(10).voltage()}},
		{"Transforms VGATE3", fields{VGATE3: 10}, CPRUReport{VGATE3: gate(10).voltage()}},
		{"Transforms VSUBS3", fields{VSUBS3: 10}, CPRUReport{VSUBS3: subs(10).voltage()}},
		{"Transforms VRD3", fields{VRD3: 10}, CPRUReport{VRD3: rd(10).voltage()}},
		{"Transforms VOD3", fields{VOD3: 10}, CPRUReport{VOD3: od(10).voltage()}},
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
