package aez

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	ccdIPD := CCDImagePackData{}
	ccdI := CCDImage{PackData: &ccdIPD, BadColumns: []uint16{42, 6, 7}}
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

func TestCCDImage_MarshalJSON(t *testing.T) {
	ccd := &CCDImage{PackData: &CCDImagePackData{}}
	got, err := ccd.MarshalJSON()
	if err != nil {
		t.Errorf("CCDImage.MarshalJSON() error = %v", err)
		return
	}
	var js map[string]interface{}
	if json.Unmarshal(got, &js) != nil {
		t.Errorf("DataRecord.MarshalJSON() = %v, not a valid json", string(got))
	}
}

func TestNewCCDImage(t *testing.T) {
	packData := CCDImagePackData{NBC: 2}
	trailing := []byte{0xff, 0xff, 0x00, 0x00, 0xcc, 0xcc}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, packData)
	data := append(buf.Bytes(), trailing...)

	tests := []struct {
		name     string
		truncate int
		want     *CCDImage
		wantErr  bool
	}{
		{
			"Returns expected",
			0,
			&CCDImage{PackData: &packData, BadColumns: []uint16{0xffff, 0x0000}},
			false,
		},
		{
			"Not enough bad columns",
			4,
			nil,
			true,
		},
		{
			"Not enough for ccd",
			24,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		reader := bytes.NewReader(data[0 : len(data)-tt.truncate])

		got, err := NewCCDImage(reader)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewCCDImage() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewCCDImage() = %v, want %v", got, tt.want)
		}
	}
}
