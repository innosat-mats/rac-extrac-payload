package ccsds

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestUnsegmentedTimeNanoSeconds(t *testing.T) {
	var tests = []struct {
		coarse uint32
		fine   uint16
		want   int64
	}{
		{0, 0, 0000000000},
		{42, 0, 42000000000},
		{0, 0x8000, 1000000000},
		{0, 0b0100000000000000, 500000000},
		{0, 0x8000 >> 8, int64(math.Round(math.Pow(2, -8) * math.Pow10(9)))},
		{0, (0x8000 >> 1) | (0x8000 >> 2), 750000000},
		{42, 0x8000 >> 3, 42000000000 + int64(math.Round(math.Pow(2, -3)*math.Pow10(9)))},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("coarse=%d,fine=%d", tt.coarse, tt.fine), func(t *testing.T) {
			got := UnsegmentedTimeNanoseconds(tt.coarse, tt.fine)
			if got != tt.want {
				t.Errorf("UnsegmentedTime(%d, %d) = %d, want %d", tt.coarse, tt.fine, got, tt.want)
			}
		})
	}
}

func TestUnsegmentedTimeDate(t *testing.T) {
	type args struct {
		coarseTime uint32
		fineTime   uint16
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"Returns Epoch/TAI", args{0, 0}, TAI},
		{"Returns expected time", args{10, 2}, TAI.Add(time.Second*10 + time.Millisecond*500)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnsegmentedTimeDate(tt.args.coarseTime, tt.args.fineTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnsegmentedTimeDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
