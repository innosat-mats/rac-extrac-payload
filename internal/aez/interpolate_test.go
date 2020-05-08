package aez

import (
	"testing"
)

func Test_getXIndex(t *testing.T) {
	type args struct {
		res         float64
		resistances []float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"getXIndex gets correct index",
			args{3e5, htrResistances[:]},
			4,
		},
		{
			"getXIndex returns index 0 if too large resistance",
			args{1e9, htrResistances[:]},
			0,
		},
		{
			"getXIndex returns last index if too small resistance",
			args{0, htrResistances[:]},
			42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getXIndex(tt.args.res, tt.args.resistances); got != tt.want {
				t.Errorf("getXIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interpolate(t *testing.T) {
	type args struct {
		r          [2]float64
		t          [2]float64
		resistance float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"interpolate is correct",
			args{[2]float64{2, 1}, [2]float64{3, 4}, 1.5},
			3.5,
		},
		{
			"interpolate is correct even outside range",
			args{[2]float64{2, 1}, [2]float64{3, 4}, 0.5},
			4.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interpolate(tt.args.r, tt.args.t, tt.args.resistance); got != tt.want {
				t.Errorf("interpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Interpolate(t *testing.T) {
	type args struct {
		resistance   float64
		resistances  []float64
		temperatures []float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Interpolate returns correct temperature",
			args{1e5, htrResistances[:], htrTemperatures[:]},
			-20.43954395439544,
			false,
		},
		{
			"Interpolate returns lowest temperature if resistance too large",
			args{10e5, htrResistances[:], htrTemperatures[:]},
			-55,
			true,
		},
		{
			"Interpolate returns highest temperature if resistance too small",
			args{0, htrResistances[:], htrTemperatures[:]},
			155,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Interpolate(tt.args.resistance, tt.args.resistances, tt.args.temperatures)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interpolate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Interpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}
