package aez

import (
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
)

func TestTCAcceptSuccess_CSVSpecifications(t *testing.T) {
	tcv := TCAcceptSuccessData{}
	want := []string{"AEZ", Specification}
	if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptSuccessData.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestTCAcceptFailure_CSVSpecifications(t *testing.T) {
	tcv := TCAcceptFailureData{}
	want := []string{"AEZ", Specification}
	if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptFailureData.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestTCExecSuccess_CSVSpecifications(t *testing.T) {
	tcv := TCExecSuccessData{}
	want := []string{"AEZ", Specification}
	if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecSuccessData.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestTCExecFailure_CSVSpecifications(t *testing.T) {
	tcv := TCExecFailureData{}
	want := []string{"AEZ", Specification}
	if got := tcv.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecFailureData.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestTCAcceptSuccess_CSVHeaders(t *testing.T) {
	tcv := TCAcceptSuccessData{}
	want := []string{"TCV", "TCPID", "PSC", "ErrorCode"}
	if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptSuccessData.CSVHeaders() = %v, want %v", got, want)
	}
}

func TestTCAcceptFailure_CSVHeaders(t *testing.T) {
	tcv := TCAcceptFailureData{}
	want := []string{"TCV", "TCPID", "PSC", "ErrorCode"}
	if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptFailureData.CSVHeaders() = %v, want %v", got, want)
	}
}

func TestTCExecSuccess_CSVHeaders(t *testing.T) {
	tcv := TCExecSuccessData{}
	want := []string{"TCV", "TCPID", "PSC", "ErrorCode"}
	if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecSuccessData.CSVHeaders() = %v, want %v", got, want)
	}
}

func TestTCExecFailure_CSVHeaders(t *testing.T) {
	tcv := TCExecFailureData{}
	want := []string{"TCV", "TCPID", "PSC", "ErrorCode"}
	if got := tcv.CSVHeaders(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecFailureData.CSVHeaders() = %v, want %v", got, want)
	}
}

func TestTCAcceptSuccess_CSVRow(t *testing.T) {
	tcv := TCAcceptSuccessData{1, 2}
	want := []string{"Accept", "1", "2", ""}
	if got := tcv.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptSuccessData.CSVRow() = %v, want %v", got, want)
	}
}

func TestTCAcceptFailure_CSVRow(t *testing.T) {
	tcv := TCAcceptFailureData{1, 2, 3}
	want := []string{"Accept", "1", "2", "3"}
	if got := tcv.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCAcceptFailureData.CSVRow() = %v, want %v", got, want)
	}
}

func TestTCExecSuccess_CSVRow(t *testing.T) {
	tcv := TCExecSuccessData{1, 2}
	want := []string{"Exec", "1", "2", ""}
	if got := tcv.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecSuccessData.CSVRow() = %v, want %v", got, want)
	}
}

func TestTCExecFailure_CSVRow(t *testing.T) {
	tcv := TCExecFailureData{1, 2, 3}
	want := []string{"Exec", "1", "2", "3"}
	if got := tcv.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("TCExecFailureData.CSVRow() = %v, want %v", got, want)
	}
}

func TestTCAcceptSuccessData_SetParquet(t *testing.T) {
	tcv := TCAcceptSuccessData{1, 2}
	want := parquetrow.ParquetRow{
		TCV:   "Accept",
		TCPID: 1,
		PSC:   2,
	}
	row := parquetrow.ParquetRow{}
	if tcv.SetParquet(&row); !reflect.DeepEqual(row, want) {
		t.Errorf("TCAcceptSuccessData.SetParquet() = %v, want %v", row, want)
	}
}

func TestTCAcceptFailureData_SetParquet(t *testing.T) {
	tcv := TCAcceptFailureData{1, 2, 3}
	want := parquetrow.ParquetRow{
		TCV:       "Accept",
		TCPID:     1,
		PSC:       2,
		ErrorCode: 3,
	}
	row := parquetrow.ParquetRow{}
	if tcv.SetParquet(&row); !reflect.DeepEqual(row, want) {
		t.Errorf("TCAcceptFailureData.SetParquet() = %v, want %v", row, want)
	}
}

func TestTCExecSuccessData_SetParquet(t *testing.T) {
	tcv := TCExecSuccessData{1, 2}
	want := parquetrow.ParquetRow{
		TCV:   "Exec",
		TCPID: 1,
		PSC:   2,
	}
	row := parquetrow.ParquetRow{}
	if tcv.SetParquet(&row); !reflect.DeepEqual(row, want) {
		t.Errorf("TCExecSuccessData.SetParquet() = %v, want %v", row, want)
	}
}

func TestTCExecFailureData_SetParquet(t *testing.T) {
	tcv := TCExecFailureData{1, 2, 3}
	want := parquetrow.ParquetRow{
		TCV:       "Exec",
		TCPID:     1,
		PSC:       2,
		ErrorCode: 3,
	}
	row := parquetrow.ParquetRow{}
	if tcv.SetParquet(&row); !reflect.DeepEqual(row, want) {
		t.Errorf("TCExecFailureData.SetParquet() = %v, want %v", row, want)
	}
}
