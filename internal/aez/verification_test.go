package aez

import (
	"reflect"
	"testing"
)

func TestTCAcceptSuccess_CSVSpecifications(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
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
			tcv := TCAcceptSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptSuccess.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCAcceptFailure_CSVSpecifications(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
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
			tcv := TCAcceptFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptFailure.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecSuccess_CSVSpecifications(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
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
			tcv := TCExecSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecSuccess.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecFailure_CSVSpecifications(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
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
			tcv := TCExecFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecFailure.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCAcceptSuccess_CSVHeaders(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates headers",
			fields{},
			[]string{"TCPID", "PSC"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCAcceptSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptSuccess.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCAcceptFailure_CSVHeaders(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates headers",
			fields{},
			[]string{"TCPID", "PSC", "ErrorCode"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCAcceptFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptFailure.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecSuccess_CSVHeaders(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates headers",
			fields{},
			[]string{"TCPID", "PSC"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCExecSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecSuccess.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecFailure_CSVHeaders(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates headers",
			fields{},
			[]string{"TCPID", "PSC", "ErrorCode"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCExecFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecFailure.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCAcceptSuccess_CSVRow(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{TCPID: 1, PSC: 2},
			[]string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCAcceptSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptSuccess.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCAcceptFailure_CSVRow(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{TCPID: 1, PSC: 2, ErrorCode: 3},
			[]string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCAcceptFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCAcceptFailure.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecSuccess_CSVRow(t *testing.T) {
	type fields struct {
		TCPID uint16
		PSC   uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{TCPID: 1, PSC: 2},
			[]string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCExecSuccess{
				TCPID: tt.fields.TCPID,
				PSC:   tt.fields.PSC,
			}
			if got := tcv.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecSuccess.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCExecFailure_CSVRow(t *testing.T) {
	type fields struct {
		TCPID     uint16
		PSC       uint16
		ErrorCode uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates a row",
			fields{TCPID: 1, PSC: 2, ErrorCode: 3},
			[]string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcv := TCExecFailure{
				TCPID:     tt.fields.TCPID,
				PSC:       tt.fields.PSC,
				ErrorCode: tt.fields.ErrorCode,
			}
			if got := tcv.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCExecFailure.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
