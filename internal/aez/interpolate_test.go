package aez

import (
	"testing"
)

func Test_getResistanceIndex(t *testing.T) {
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
			"getResistanceIndex gets correct index",
			args{3e5, htrResistances[:]},
			4,
		},
		{
			"getResistanceIndex returns index 0 if too large resistance",
			args{1e9, htrResistances[:]},
			0,
		},
		{
			"getResistanceIndex returns last index if too small resistance",
			args{0, htrResistances[:]},
			42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResistanceIndex(tt.args.res, tt.args.resistances); got != tt.want {
				t.Errorf("getResistanceIndex() = %v, want %v", got, tt.want)
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

func Test_interpolateTemperature(t *testing.T) {
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
			"interpolateTemperature returns correct temperature",
			args{1e5, htrResistances[:], htrTemperatures[:]},
			-20.43954395439544,
			false,
		},
		{
			"interpolateTemperature returns lowest temperature if resistance too large",
			args{10e5, htrResistances[:], htrTemperatures[:]},
			-55,
			true,
		},
		{
			"interpolateTemperature returns highest temperature if resistance too small",
			args{0, htrResistances[:], htrTemperatures[:]},
			155,
			true,
		},
		{
			"interpolateTemperature gets angry if resistances and temperatures of different size",
			args{0, htrResistances[:], pwrTemperatures[:]},
			-273.15,
			true,
		},
		{
			"interpolateTemperature gets angry if resistances or temperatures too short",
			args{0, []float64{1}, []float64{0}},
			-273.15,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interpolateTemperature(tt.args.resistance, tt.args.resistances, tt.args.temperatures)
			if (err != nil) != tt.wantErr {
				t.Errorf("interpolateTemperature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("interpolateTemperature() = %v, want %v", got, tt.want)
			}
		})
	}
}
