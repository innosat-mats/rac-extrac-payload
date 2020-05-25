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
	pwrt10 := pwrt(10)
	temperature, _ := pwrt10.temperature()
	pwrp32v10 := pwrp32v(10)
	pwrp32c10 := pwrp32c(10)
	pwrp16v10 := pwrp16v(10)
	pwrp16c10 := pwrp16c(10)
	pwrm16v10 := pwrm16v(10)
	pwrm16c10 := pwrm16c(10)
	pwrp3v310 := pwrp3v3(10)
	pwrp3c310 := pwrp3c3(10)
	tests := []struct {
		name   string
		fields fields
		field  string
		want   float64
	}{
		{"PWRT is temperature", fields{PWRT: 10}, "PWRT", temperature},
		{"PWRP32V is voltage", fields{PWRP32V: 10}, "PWRP32V", pwrp32v10.voltage()},
		{"PWRP32C is current", fields{PWRP32C: 10}, "PWRP32C", pwrp32c10.current()},
		{"PWRP16V is voltage", fields{PWRP16V: 10}, "PWRP16V", pwrp16v10.voltage()},
		{"PWRP16C is current", fields{PWRP16C: 10}, "PWRP16C", pwrp16c10.current()},
		{"PWRM16V is voltage", fields{PWRM16V: 10}, "PWRM16V", pwrm16v10.voltage()},
		{"PWRM16C is current", fields{PWRM16C: 10}, "PWRM16C", pwrm16c10.current()},
		{"PWRP3V3 is voltage", fields{PWRP3V3: 10}, "PWRP3V3", pwrp3v310.voltage()},
		{"PWRP3C3 is current", fields{PWRP3C3: 10}, "PWRP3C3", pwrp3c310.current()},
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
	tests := []struct {
		name string
		want []string
	}{
		{"Genereates spec", []string{"AEZ", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwr := PWR{}
			if got := pwr.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PWR.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPWR_CSVHeaders(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			"Generates headers",
			[]string{
				"PWRT",
				"PWRP32V", "PWRP32C",
				"PWRP16V", "PWRP16C",
				"PWRM16V", "PWRM16C",
				"PWRP3V3", "PWRP3C3",
				"WARNINGS",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwr := PWR{}
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
	pwrp32v2 := pwrp32v(2)
	pwrp32c3 := pwrp32c(3)
	pwrp16v4 := pwrp16v(4)
	pwrp16c5 := pwrp16c(5)
	pwrm16v6 := pwrm16v(6)
	pwrm16c7 := pwrm16c(7)
	pwrp3v38 := pwrp3v3(8)
	pwrp3c39 := pwrp3c3(9)
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
				"-55",
				fmt.Sprintf("%v", pwrp32v2.voltage()),
				fmt.Sprintf("%v", pwrp32c3.current()),
				fmt.Sprintf("%v", pwrp16v4.voltage()),
				fmt.Sprintf("%v", pwrp16c5.current()),
				fmt.Sprintf("%v", pwrm16v6.voltage()),
				fmt.Sprintf("%v", pwrm16c7.current()),
				fmt.Sprintf("%v", pwrp3v38.voltage()),
				fmt.Sprintf("%v", pwrp3c39.current()),
				"PWRT: 5.4044e+06 is too large for interpolator. Returning value for maximum.",
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
