package aez

import (
	"fmt"
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
	temperature, _ := htr(10).temperature()
	tests := []struct {
		name   string
		fields fields
		field  string
		want   float64
	}{
		{"HTR1A is temperature", fields{HTR1A: 10}, "HTR1A", temperature},
		{"HTR1B is temperature", fields{HTR1B: 10}, "HTR1B", temperature},
		{"HTR1OD is voltage", fields{HTR1OD: 10}, "HTR1OD", htr(10).voltage()},
		{"HTR2A is temperature", fields{HTR2A: 10}, "HTR2A", temperature},
		{"HTR2B is temperature", fields{HTR2B: 10}, "HTR2B", temperature},
		{"HTR2OD is voltage", fields{HTR2OD: 10}, "HTR2OD", htr(10).voltage()},
		{"HTR7A is temperature", fields{HTR7A: 10}, "HTR7A", temperature},
		{"HTR7B is temperature", fields{HTR7B: 10}, "HTR7B", temperature},
		{"HTR7OD is voltage", fields{HTR7OD: 10}, "HTR7OD", htr(10).voltage()},
		{"HTR8A is temperature", fields{HTR8A: 10}, "HTR8A", temperature},
		{"HTR8B is temperature", fields{HTR8B: 10}, "HTR8B", temperature},
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

func TestHTR_CSVSpecifications(t *testing.T) {
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
		want   []string
	}{
		{"Genereates spec", fields{}, []string{"AEZ", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htr := HTR{
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
			if got := htr.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTR.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTR_CSVHeaders(t *testing.T) {
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
		want   []string
	}{
		{
			"Generates expected headers",
			fields{},
			[]string{
				"HTR1A", "HTR1B", "HTR1OD",
				"HTR2A", "HTR2B", "HTR2OD",
				"HTR7A", "HTR7B", "HTR7OD",
				"HTR8A", "HTR8B", "HTR8OD",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htr := HTR{
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
			if got := htr.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTR.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTR_CSVRow(t *testing.T) {
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
		want   []string
	}{
		{
			"Generates a data row",
			fields{
				HTR1A: 1, HTR1B: 2, HTR1OD: 3,
				HTR2A: 4, HTR2B: 5, HTR2OD: 6,
				HTR7A: 7, HTR7B: 8, HTR7OD: 9,
				HTR8A: 10, HTR8B: 11, HTR8OD: 12,
			},
			[]string{
				fmt.Sprintf("%v", htr(1).voltage()),
				fmt.Sprintf("%v", htr(2).resistance()),
				fmt.Sprintf("%v", htr(3).voltage()),
				fmt.Sprintf("%v", htr(4).voltage()),
				fmt.Sprintf("%v", htr(5).resistance()),
				fmt.Sprintf("%v", htr(6).voltage()),
				fmt.Sprintf("%v", htr(7).voltage()),
				fmt.Sprintf("%v", htr(8).resistance()),
				fmt.Sprintf("%v", htr(9).voltage()),
				fmt.Sprintf("%v", htr(10).voltage()),
				fmt.Sprintf("%v", htr(11).resistance()),
				fmt.Sprintf("%v", htr(12).voltage()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htr := HTR{
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
			if got := htr.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTR.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
