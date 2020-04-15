package ramses

import (
	"reflect"
	"testing"
	"time"
)

func TestRamses_Created(t *testing.T) {
	tests := []struct {
		name string
		r    *Ramses
		want time.Time
	}{
		{
			"zero time",
			&Ramses{},
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"1 day",
			&Ramses{0, 0, 0, 0, 0, 0, 1},
			time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			"1.5 days",
			&Ramses{0, 0, 0, 0, 0, 43200000, 1},
			time.Date(2000, 1, 2, 12, 0, 0, 0, time.UTC),
		}, {
			"-1.5 days",
			&Ramses{0, 0, 0, 0, 0, 43200000, -2},
			time.Date(1999, 12, 30, 12, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Created(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ramses.Created() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRamses_Valid(t *testing.T) {
	tests := []struct {
		name string
		r    *Ramses
		want bool
	}{
		{
			"valid",
			&Ramses{0xeb90, 0, 0, 0, 0, 0, 0},
			true,
		}, {
			"invalid",
			&Ramses{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Valid(); got != tt.want {
				t.Errorf("Ramses.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
