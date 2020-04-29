package aez

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPWR_Report(t *testing.T) {
	type fields struct {
		PWRT    pwrt
		PWRP32V pwrp32v
		PWRP32C pwrp32c
		PWRP16V pwrp16v
		PWRP16C pwrp16c
		PWRM16V pwrm16v
		PWRM16C pwrm16c
		PWRP3V3 pwrp3v3
		PWRP3C3 pwrp3c3
	}
	tests := []struct {
		name   string
		fields fields
		field  string
		want   float64
	}{
		{"PWRT is resistance", fields{PWRT: 10}, "PWRT", pwrt(10).resistance()},
		{"PWRP32V is voltage", fields{PWRP32V: 10}, "PWRP32V", pwrp32v(10).voltage()},
		{"PWRP32C is current", fields{PWRP32C: 10}, "PWRP32C", pwrp32c(10).current()},
		{"PWRP16V is voltage", fields{PWRP16V: 10}, "PWRP16V", pwrp16v(10).voltage()},
		{"PWRP16C is current", fields{PWRP16C: 10}, "PWRP16C", pwrp16c(10).current()},
		{"PWRM16V is voltage", fields{PWRM16V: 10}, "PWRM16V", pwrm16v(10).voltage()},
		{"PWRM16C is current", fields{PWRM16C: 10}, "PWRM16C", pwrm16c(10).current()},
		{"PWRP3V3 is voltage", fields{PWRP3V3: 10}, "PWRP3V3", pwrp3v3(10).voltage()},
		{"PWRP3C3 is current", fields{PWRP3C3: 10}, "PWRP3C3", pwrp3c3(10).current()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwr := &PWR{
				PWRT:    tt.fields.PWRT,
				PWRP32V: tt.fields.PWRP32V,
				PWRP32C: tt.fields.PWRP32C,
				PWRP16V: tt.fields.PWRP16V,
				PWRP16C: tt.fields.PWRP16C,
				PWRM16V: tt.fields.PWRM16V,
				PWRM16C: tt.fields.PWRM16C,
				PWRP3V3: tt.fields.PWRP3V3,
				PWRP3C3: tt.fields.PWRP3C3,
			}
			report := pwr.Report()
			value := reflect.ValueOf(report)
			field := reflect.Indirect(value).FieldByName(tt.field)
			if got := field.Float(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PWR.Report() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPWR_CSVSpecifications(t *testing.T) {
	type fields struct {
		PWRT    pwrt
		PWRP32V pwrp32v
		PWRP32C pwrp32c
		PWRP16V pwrp16v
		PWRP16C pwrp16c
		PWRM16V pwrm16v
		PWRM16C pwrm16c
		PWRP3V3 pwrp3v3
		PWRP3C3 pwrp3c3
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
			pwr := PWR{
				PWRT:    tt.fields.PWRT,
				PWRP32V: tt.fields.PWRP32V,
				PWRP32C: tt.fields.PWRP32C,
				PWRP16V: tt.fields.PWRP16V,
				PWRP16C: tt.fields.PWRP16C,
				PWRM16V: tt.fields.PWRM16V,
				PWRM16C: tt.fields.PWRM16C,
				PWRP3V3: tt.fields.PWRP3V3,
				PWRP3C3: tt.fields.PWRP3C3,
			}
			if got := pwr.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PWR.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPWR_CSVHeaders(t *testing.T) {
	type fields struct {
		PWRT    pwrt
		PWRP32V pwrp32v
		PWRP32C pwrp32c
		PWRP16V pwrp16v
		PWRP16C pwrp16c
		PWRM16V pwrm16v
		PWRM16C pwrm16c
		PWRP3V3 pwrp3v3
		PWRP3C3 pwrp3c3
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
				"PWRT",
				"PWRP32V", "PWRP32C",
				"PWRP16V", "PWRP16C",
				"PWRM16V", "PWRM16C",
				"PWRP3V3", "PWRP3C3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwr := PWR{
				PWRT:    tt.fields.PWRT,
				PWRP32V: tt.fields.PWRP32V,
				PWRP32C: tt.fields.PWRP32C,
				PWRP16V: tt.fields.PWRP16V,
				PWRP16C: tt.fields.PWRP16C,
				PWRM16V: tt.fields.PWRM16V,
				PWRM16C: tt.fields.PWRM16C,
				PWRP3V3: tt.fields.PWRP3V3,
				PWRP3C3: tt.fields.PWRP3C3,
			}
			if got := pwr.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PWR.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPWR_CSVRow(t *testing.T) {
	type fields struct {
		PWRT    pwrt
		PWRP32V pwrp32v
		PWRP32C pwrp32c
		PWRP16V pwrp16v
		PWRP16C pwrp16c
		PWRM16V pwrm16v
		PWRM16C pwrm16c
		PWRP3V3 pwrp3v3
		PWRP3C3 pwrp3c3
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates data row",
			fields{
				PWRT:    pwrt(1),
				PWRP32V: pwrp32v(2),
				PWRP32C: pwrp32c(3),
				PWRP16V: pwrp16v(4),
				PWRP16C: pwrp16c(5),
				PWRM16V: pwrm16v(6),
				PWRM16C: pwrm16c(7),
				PWRP3V3: pwrp3v3(8),
				PWRP3C3: pwrp3c3(9),
			},
			[]string{
				fmt.Sprintf("%v", pwrt(1).resistance()),
				fmt.Sprintf("%v", pwrp32v(2).voltage()),
				fmt.Sprintf("%v", pwrp32c(3).current()),
				fmt.Sprintf("%v", pwrp16v(4).voltage()),
				fmt.Sprintf("%v", pwrp16c(5).current()),
				fmt.Sprintf("%v", pwrm16v(6).voltage()),
				fmt.Sprintf("%v", pwrm16c(7).current()),
				fmt.Sprintf("%v", pwrp3v3(8).voltage()),
				fmt.Sprintf("%v", pwrp3c3(9).current()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwr := PWR{
				PWRT:    tt.fields.PWRT,
				PWRP32V: tt.fields.PWRP32V,
				PWRP32C: tt.fields.PWRP32C,
				PWRP16V: tt.fields.PWRP16V,
				PWRP16C: tt.fields.PWRP16C,
				PWRM16V: tt.fields.PWRM16V,
				PWRM16C: tt.fields.PWRM16C,
				PWRP3V3: tt.fields.PWRP3V3,
				PWRP3C3: tt.fields.PWRP3C3,
			}
			if got := pwr.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PWR.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
