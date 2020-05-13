package aez

import (
	"reflect"
	"testing"
)

func TestCCDImage_CSVSpecifications(t *testing.T) {
	ccd := CCDImage{}
	want := []string{"Specification", Specification}
	if got := ccd.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("CCDImage.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestCCDImage_CSVHeaders_AddsBC(t *testing.T) {
	ccdI := CCDImage{}
	ccdIPD := CCDImagePackData{}
	headersI := ccdI.CSVHeaders()
	headersIPD := ccdIPD.CSVHeaders()
	wantBC := "BC"
	for i, header := range headersI {
		if i < len(headersIPD) {
			if header != headersIPD[i] {
				t.Errorf("%v: got %v, want %v", i, header, headersIPD[i])
			}
		} else if i == len(headersIPD) {
			if header != wantBC {
				t.Errorf("%v: got %v, want %v", i, header, wantBC)
			}
		} else {
			t.Errorf("Unexpected %vth header %v", i, header)
		}
	}
}

func TestCCDImage_CSVRow_AddsBC(t *testing.T) {
	ccdI := CCDImage{BadColumns: []uint16{42, 6, 7}}
	ccdIPD := CCDImagePackData{}
	rowI := ccdI.CSVRow()
	rowIPD := ccdIPD.CSVRow()
	wantBC := "[42 6 7]"
	for i, value := range rowI {
		if i < len(rowIPD) {
			if value != rowIPD[i] {
				t.Errorf("%v: got %v, want %v", i, value, rowIPD[i])
			}
		} else if i == len(rowIPD) {
			if value != wantBC {
				t.Errorf("%v: got %v, want %v", i, value, wantBC)
			}
		} else {
			t.Errorf("Unexpected %vth column %v", i, value)
		}
	}
}
