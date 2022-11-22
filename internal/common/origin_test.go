package common

import (
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
)

func TestOriginDescription_CSVHeaders(t *testing.T) {
	type fields struct {
		Name           string
		ProcessingDate time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Returns headers", fields{}, []string{"OriginFile", "ProcessingDate"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := OriginDescription{
				Name:           tt.fields.Name,
				ProcessingDate: tt.fields.ProcessingDate,
			}
			if got := origin.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OriginDescription.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOriginDescription_CSVRow(t *testing.T) {
	type fields struct {
		Name           string
		ProcessingDate time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Returns data row",
			fields{Name: "Sir Longtailed Tit", ProcessingDate: procDate},
			[]string{
				"Sir Longtailed Tit",
				procDate.Format(time.RFC3339),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := OriginDescription{
				Name:           tt.fields.Name,
				ProcessingDate: tt.fields.ProcessingDate,
			}
			if got := origin.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OriginDescription.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOriginDescription_SetParquet(t *testing.T) {
	origin := OriginDescription{"Name", time.Time{}}
	want := parquetrow.ParquetRow{
		OriginFile:     "Name",
		ProcessingTime: time.Time{},
	}
	row := parquetrow.ParquetRow{}
	if origin.SetParquet(&row); !reflect.DeepEqual(row, want) {
		t.Errorf("OriginDescription.SetParquet() = %v, want %v", row, want)
	}
}
