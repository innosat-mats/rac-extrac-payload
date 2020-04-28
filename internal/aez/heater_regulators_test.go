package aez

import (
	"reflect"
	"testing"
)

func TestHTR_Report(t *testing.T) {
	type fields struct {
		HTR1A  htr
		HTR1B  htr
		HTR1OD htr
		HTR2A  htr
		HTR2B  htr
		HTR2OD htr
		HTR7A  htr
		HTR7B  htr
		HTR7OD htr
		HTR8A  htr
		HTR8B  htr
		HTR8OD htr
	}
	tests := []struct {
		name   string
		fields fields
		field  string
		want   float64
	}{
		{"HTR1A is temperature", fields{HTR1A: 10}, "HTR1A", htr(10).temperature()},
		{"HTR1B is temperature", fields{HTR1B: 10}, "HTR1B", htr(10).temperature()},
		{"HTR1OD is voltage", fields{HTR1OD: 10}, "HTR1OD", htr(10).voltage()},
		{"HTR2A is temperature", fields{HTR2A: 10}, "HTR2A", htr(10).temperature()},
		{"HTR2B is temperature", fields{HTR2B: 10}, "HTR2B", htr(10).temperature()},
		{"HTR2OD is voltage", fields{HTR2OD: 10}, "HTR2OD", htr(10).voltage()},
		{"HTR7A is temperature", fields{HTR7A: 10}, "HTR7A", htr(10).temperature()},
		{"HTR7B is temperature", fields{HTR7B: 10}, "HTR7B", htr(10).temperature()},
		{"HTR7OD is voltage", fields{HTR7OD: 10}, "HTR7OD", htr(10).voltage()},
		{"HTR8A is temperature", fields{HTR8A: 10}, "HTR8A", htr(10).temperature()},
		{"HTR8B is temperature", fields{HTR8B: 10}, "HTR8B", htr(10).temperature()},
		{"HTR8OD is voltage", fields{HTR8OD: 10}, "HTR8OD", htr(10).voltage()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htr := &HTR{
				HTR1A:  tt.fields.HTR1A,
				HTR1B:  tt.fields.HTR1B,
				HTR1OD: tt.fields.HTR1OD,
				HTR2A:  tt.fields.HTR2A,
				HTR2B:  tt.fields.HTR2B,
				HTR2OD: tt.fields.HTR2OD,
				HTR7A:  tt.fields.HTR7A,
				HTR7B:  tt.fields.HTR7B,
				HTR7OD: tt.fields.HTR7OD,
				HTR8A:  tt.fields.HTR8A,
				HTR8B:  tt.fields.HTR8B,
				HTR8OD: tt.fields.HTR8OD,
			}
			report := htr.Report()
			value := reflect.ValueOf(report)
			field := reflect.Indirect(value).FieldByName(tt.field)
			if got := field.Float(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTR.Report() = %v, want %v", got, tt.want)
			}
		})
	}
}
