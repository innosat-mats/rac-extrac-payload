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
		HTR3A  htr
		HTR3B  htr
		HTR3OD htr
		HTR4A  htr
		HTR4B  htr
		HTR4OD htr
		HTR5A  htr
		HTR5B  htr
		HTR5OD htr
		HTR6A  htr
		HTR6B  htr
		HTR6OD htr
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
		{"HTR1A is voltage", fields{HTR1A: 10}, "HTR1A", htr(10).voltage()},
		{"HTR1B is resistance", fields{HTR1B: 10}, "HTR1B", htr(10).resistance()},
		{"HTR1OD is voltage", fields{HTR1OD: 10}, "HTR1OD", htr(10).voltage()},
		{"HTR2A is voltage", fields{HTR2A: 10}, "HTR2A", htr(10).voltage()},
		{"HTR2B is resistance", fields{HTR2B: 10}, "HTR2B", htr(10).resistance()},
		{"HTR2OD is voltage", fields{HTR2OD: 10}, "HTR2OD", htr(10).voltage()},
		{"HTR3A is voltage", fields{HTR3A: 10}, "HTR3A", htr(10).voltage()},
		{"HTR3B is resistance", fields{HTR3B: 10}, "HTR3B", htr(10).resistance()},
		{"HTR3OD is voltage", fields{HTR3OD: 10}, "HTR3OD", htr(10).voltage()},
		{"HTR4A is voltage", fields{HTR4A: 10}, "HTR4A", htr(10).voltage()},
		{"HTR4B is resistance", fields{HTR4B: 10}, "HTR4B", htr(10).resistance()},
		{"HTR4OD is voltage", fields{HTR4OD: 10}, "HTR4OD", htr(10).voltage()},
		{"HTR5A is voltage", fields{HTR5A: 10}, "HTR5A", htr(10).voltage()},
		{"HTR5B is resistance", fields{HTR5B: 10}, "HTR5B", htr(10).resistance()},
		{"HTR5OD is voltage", fields{HTR5OD: 10}, "HTR5OD", htr(10).voltage()},
		{"HTR6A is voltage", fields{HTR6A: 10}, "HTR6A", htr(10).voltage()},
		{"HTR6B is resistance", fields{HTR6B: 10}, "HTR6B", htr(10).resistance()},
		{"HTR6OD is voltage", fields{HTR6OD: 10}, "HTR6OD", htr(10).voltage()},
		{"HTR7A is voltage", fields{HTR7A: 10}, "HTR7A", htr(10).voltage()},
		{"HTR7B is resistance", fields{HTR7B: 10}, "HTR7B", htr(10).resistance()},
		{"HTR7OD is voltage", fields{HTR7OD: 10}, "HTR7OD", htr(10).voltage()},
		{"HTR8A is voltage", fields{HTR8A: 10}, "HTR8A", htr(10).voltage()},
		{"HTR8B is resistance", fields{HTR8B: 10}, "HTR8B", htr(10).resistance()},
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
				HTR3A:  tt.fields.HTR3A,
				HTR3B:  tt.fields.HTR3B,
				HTR3OD: tt.fields.HTR3OD,
				HTR4A:  tt.fields.HTR4A,
				HTR4B:  tt.fields.HTR4B,
				HTR4OD: tt.fields.HTR4OD,
				HTR5A:  tt.fields.HTR5A,
				HTR5B:  tt.fields.HTR5B,
				HTR5OD: tt.fields.HTR5OD,
				HTR6A:  tt.fields.HTR6A,
				HTR6B:  tt.fields.HTR6B,
				HTR6OD: tt.fields.HTR6OD,
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
