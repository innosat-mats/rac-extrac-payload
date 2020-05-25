package ramses

import (
	"bytes"
	"fmt"
	"io"
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

func TestRamses_CSVSpecifications(t *testing.T) {
	type fields struct {
		Synch  uint16
		Length uint16
		Port   uint16
		Type   uint8
		Secure uint8
		Time   uint32
		Date   int32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Creates spec", fields{}, []string{"RAMSES", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ramses := Ramses{
				Synch:  tt.fields.Synch,
				Length: tt.fields.Length,
				Port:   tt.fields.Port,
				Type:   tt.fields.Type,
				Secure: tt.fields.Secure,
				Time:   tt.fields.Time,
				Date:   tt.fields.Date,
			}
			if got := ramses.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ramses.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRamses_CSVHeaders(t *testing.T) {
	type fields struct {
		Synch  uint16
		Length uint16
		Port   uint16
		Type   uint8
		Secure uint8
		Time   uint32
		Date   int32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Creates headers", fields{}, []string{"RamsesTime"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ramses := Ramses{
				Synch:  tt.fields.Synch,
				Length: tt.fields.Length,
				Port:   tt.fields.Port,
				Type:   tt.fields.Type,
				Secure: tt.fields.Secure,
				Time:   tt.fields.Time,
				Date:   tt.fields.Date,
			}
			if got := ramses.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ramses.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRamses_CSVRow(t *testing.T) {
	type fields struct {
		Synch  uint16
		Length uint16
		Port   uint16
		Type   uint8
		Secure uint8
		Time   uint32
		Date   int32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates data row",
			fields{Date: 24, Time: 42},
			[]string{"2000-01-25T00:00:00.042Z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ramses := Ramses{
				Synch:  tt.fields.Synch,
				Length: tt.fields.Length,
				Port:   tt.fields.Port,
				Type:   tt.fields.Type,
				Secure: tt.fields.Secure,
				Time:   tt.fields.Time,
				Date:   tt.fields.Date,
			}
			if got := ramses.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ramses.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRamses_MarshalJSON(t *testing.T) {
	type fields struct {
		Synch  uint16
		Length uint16
		Port   uint16
		Type   uint8
		Secure uint8
		Time   uint32
		Date   int32
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"Marshals into expected json",
			fields{},
			[]byte(fmt.Sprintf("{\"specification\":\"%v\",\"ramsesTime\":\"2000-01-01T00:00:00Z\"}", Specification)),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ramses := &Ramses{
				Synch:  tt.fields.Synch,
				Length: tt.fields.Length,
				Port:   tt.fields.Port,
				Type:   tt.fields.Type,
				Secure: tt.fields.Secure,
				Time:   tt.fields.Time,
				Date:   tt.fields.Date,
			}
			got, err := ramses.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ramses.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ramses.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRamses_Nanoseconds(t *testing.T) {
	type fields struct {
		Synch  uint16
		Length uint16
		Port   uint16
		Type   uint8
		Secure uint8
		Time   uint32
		Date   int32
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"Returns 0",
			fields{},
			0,
		},
		{
			"Returns expected value",
			fields{Time: 42, Date: 10},
			10*3600*24*1000000000 + 42*1000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ramses := &Ramses{
				Synch:  tt.fields.Synch,
				Length: tt.fields.Length,
				Port:   tt.fields.Port,
				Type:   tt.fields.Type,
				Secure: tt.fields.Secure,
				Time:   tt.fields.Time,
				Date:   tt.fields.Date,
			}
			if got := ramses.Nanoseconds(); got != tt.want {
				t.Errorf("Ramses.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRamses(t *testing.T) {
	tests := []struct {
		name          string
		buf           []byte
		wantErr       bool
		wantEOF       bool
		ramsesOutcome Ramses
	}{
		{
			"Returns EOF if no bytes in buffer",
			[]byte{},
			true,
			true,
			Ramses{},
		},
		{
			"Returns non-EOF err on shorter than needed buffer",
			[]byte{0xaa, 0xab},
			true,
			false,
			Ramses{},
		},
		{
			"Reads into ramses",
			append([]byte{0x90, 0xeb, 0x48, 0x00}, make([]byte, 12)...),
			false,
			false,
			Ramses{Synch: 0xeb90, Length: 0x0048},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.buf)
			ramses, err := NewRamses(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ramses.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == io.EOF) != tt.wantEOF {
				t.Errorf("Ramses.Read() error = %v, wantEOF %v", err, tt.wantEOF)
			}
			if !tt.wantErr && !reflect.DeepEqual(*ramses, tt.ramsesOutcome) {
				t.Errorf("Ramses.Read() => %+v, want %+v", ramses, tt.ramsesOutcome)
			}
		})
	}
}
