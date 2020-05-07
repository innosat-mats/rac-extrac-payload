package aez

import (
	"reflect"
	"testing"
)

func TestCCDImage_CSVSpecifications(t *testing.T) {
	type fields struct {
		PackData   CCDImagePackData
		BadColumns []uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Returns spec array", fields{}, []string{"Specification", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccd := CCDImage{
				PackData:   tt.fields.PackData,
				BadColumns: tt.fields.BadColumns,
			}
			if got := ccd.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImage.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImage_CSVHeaders(t *testing.T) {
	type fields struct {
		PackData   CCDImagePackData
		BadColumns []uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Retruns no headers (for now)", fields{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccd := CCDImage{
				PackData:   tt.fields.PackData,
				BadColumns: tt.fields.BadColumns,
			}
			if got := ccd.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImage.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImage_CSVRow(t *testing.T) {
	type fields struct {
		PackData   CCDImagePackData
		BadColumns []uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Returns no data row (for now)", fields{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccd := CCDImage{
				PackData:   tt.fields.PackData,
				BadColumns: tt.fields.BadColumns,
			}
			if got := ccd.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImage.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
