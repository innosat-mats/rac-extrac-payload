package aez

import (
	"reflect"
	"testing"
)

func TestCPRU_Report(t *testing.T) {
	type fields struct {
		VGATE0 uint16
		VSUBS0 uint16
		VRD0   uint16
		VOD0   uint16
		VGATE1 uint16
		VSUBS1 uint16
		VRD1   uint16
		VOD1   uint16
		VGATE2 uint16
		VSUBS2 uint16
		VRD2   uint16
		VOD2   uint16
		VGATE3 uint16
		VSUBS3 uint16
		VRD3   uint16
		VOD3   uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   CPRUReport
	}{
		{"Transforms VGATE0", fields{VGATE0: 10}, CPRUReport{VGATE0: gateVoltage(10)}},
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
