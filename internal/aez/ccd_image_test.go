package aez

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCCDImage_CSVSpecifications(t *testing.T) {
	ccd := CCDImage{}
	want := []string{"AEZ", Specification}
	if got := ccd.CSVSpecifications(); !reflect.DeepEqual(got, want) {
		t.Errorf("CCDImage.CSVSpecifications() = %v, want %v", got, want)
	}
}

func TestCCDImage_CSVHeaders_AddsOwn(t *testing.T) {
	ccdI := CCDImage{}
	ccdIPD := CCDImagePackData{}
	headersI := ccdI.CSVHeaders()
	want := append(ccdIPD.CSVHeaders(), "BC", "Image File Name")

	for i, header := range headersI {
		if i < len(want) {
			if header != want[i] {
				t.Errorf("%v: got %v, want %v", i, header, want[i])
			}
		} else {
			t.Errorf("Unexpected %vth header %v", i, header)
		}
	}
	if len(headersI) < len(want) {
		t.Errorf(
			"Got %v headers, want %v (missing %v)",
			len(headersI),
			len(want),
			want[len(headersI):],
		)
	}
}

func TestCCDImage_CSVRow_AddsOwn(t *testing.T) {
	ccdIPD := CCDImagePackData{}
	ccdI := CCDImage{PackData: &ccdIPD, BadColumns: []uint16{42, 6, 7}, ImageFileName: "my_🖼️.png"}
	rowI := ccdI.CSVRow()
	want := append(ccdIPD.CSVRow(), "[42 6 7]", "my_🖼️.png")
	for i, value := range rowI {
		if i < len(want) {
			if value != want[i] {
				t.Errorf("%v: got %v, want %v", i, value, want[i])
			}
		} else {
			t.Errorf("Unexpected %vth column %v", i, value)
		}
	}
	if len(rowI) < len(want) {
		t.Errorf(
			"Got %v headers, want %v (missing %v)",
			len(rowI),
			len(want),
			want[len(rowI):],
		)
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
		name       string
		truncate   int
		originName string
		rid        RID
		want       *CCDImage
		wantErr    bool
	}{
		{
			"Returns expected",
			0,
			"my_rac.rac",
			CCD1,
			&CCDImage{
				PackData:      &packData,
				BadColumns:    []uint16{0xffff, 0x0000},
				ImageFileName: "my_rac_0_1.png",
			},
			false,
		},
		{
			"Not enough bad columns",
			4,
			"my_rac.rac",
			CCD1,
			nil,
			true,
		},
		{
			"Not enough for ccd",
			24,
			"my_rac.rac",
			CCD1,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		reader := bytes.NewReader(data[0 : len(data)-tt.truncate])

		got, err := NewCCDImage(reader, tt.originName, tt.rid)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewCCDImage() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewCCDImage() = %v, want %v", got, tt.want)
		}
	}
}

func TestCCDImage_FullImageName(t *testing.T) {
	tests := []struct {
		name          string
		imageFileName string
		prefix        string
		want          string
	}{
		{"no prefix just filename", "test.png", "", "test.png"},
		{"filename with prefix", "😓️.png", "🌞️", filepath.Join("🌞️", "😓️.png")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccd := &CCDImage{
				ImageFileName: tt.imageFileName,
			}
			if got := ccd.FullImageName(tt.prefix); got != tt.want {
				t.Errorf("CCDImage.FullImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
