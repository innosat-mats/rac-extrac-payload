package aez

import (
	"math"
	"testing"
)

func TestCPRU_Voltage(t *testing.T) {
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
	type args struct {
		field CPRUField
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Calculates voltage for VGATE0",
			fields{VGATE0: 10},
			args{VGATE0},
			2.5 / (math.Pow(2, 12) - 1) * 10 * 10,
			false,
		},
		{
			"Calculates voltage for VGATE1",
			fields{VGATE1: 10},
			args{VGATE1},
			2.5 / (math.Pow(2, 12) - 1) * 10 * 10,
			false,
		},
		{
			"Calculates voltage for VGATE2",
			fields{VGATE2: 10},
			args{VGATE2},
			2.5 / (math.Pow(2, 12) - 1) * 10 * 10,
			false,
		},
		{
			"Calculates voltage for VGATE3",
			fields{VGATE3: 10},
			args{VGATE3},
			2.5 / (math.Pow(2, 12) - 1) * 10 * 10,
			false,
		},
		{
			"Calculates voltage for VSUBS0",
			fields{VSUBS0: 1},
			args{VSUBS0},
			2.5 / (math.Pow(2, 12) - 1) * 1 * 11 / 1.5,
			false,
		},
		{
			"Calculates voltage for VSUBS1",
			fields{VSUBS1: 1},
			args{VSUBS1},
			2.5 / (math.Pow(2, 12) - 1) * 1 * 11 / 1.5,
			false,
		},
		{
			"Calculates voltage for VSUBS2",
			fields{VSUBS2: 1},
			args{VSUBS2},
			2.5 / (math.Pow(2, 12) - 1) * 1 * 11 / 1.5,
			false,
		},
		{
			"Calculates voltage for VSUBS3",
			fields{VSUBS3: 1},
			args{VSUBS3},
			2.5 / (math.Pow(2, 12) - 1) * 1 * 11 / 1.5,
			false,
		},
		{
			"Calculates voltage for VRD0",
			fields{VRD0: 2},
			args{VRD0},
			2.5 / (math.Pow(2, 12) - 1) * 2 * 17 / 1.5,
			false,
		},
		{
			"Calculates voltage for VRD1",
			fields{VRD1: 2},
			args{VRD1},
			2.5 / (math.Pow(2, 12) - 1) * 2 * 17 / 1.5,
			false,
		},
		{
			"Calculates voltage for VRD2",
			fields{VRD2: 2},
			args{VRD2},
			2.5 / (math.Pow(2, 12) - 1) * 2 * 17 / 1.5,
			false,
		},
		{
			"Calculates voltage for VRD3",
			fields{VRD3: 2},
			args{VRD3},
			2.5 / (math.Pow(2, 12) - 1) * 2 * 17 / 1.5,
			false,
		},
		{
			"Calculates voltage for VOD0",
			fields{VOD0: 3},
			args{VOD0},
			2.5 / (math.Pow(2, 12) - 1) * 3 * 32 / 1.5,
			false,
		},
		{
			"Calculates voltage for VOD1",
			fields{VOD1: 3},
			args{VOD1},
			2.5 / (math.Pow(2, 12) - 1) * 3 * 32 / 1.5,
			false,
		},
		{
			"Calculates voltage for VOD2",
			fields{VOD2: 3},
			args{VOD2},
			2.5 / (math.Pow(2, 12) - 1) * 3 * 32 / 1.5,
			false,
		},
		{
			"Calculates voltage for VOD3",
			fields{VOD3: 3},
			args{VOD3},
			2.5 / (math.Pow(2, 12) - 1) * 3 * 32 / 1.5,
			false,
		},
		{
			"Panics for unknown CPRUField",
			fields{},
			args{-10},
			0,
			true,
		},
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
			got, err := cpru.Voltage(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("CPRU.Voltage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CPRU.Voltage() = %v, want %v", got, tt.want)
			}
		})
	}
}
