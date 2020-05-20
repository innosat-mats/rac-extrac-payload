package timeseries

import (
	"bytes"
	"strings"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func TestNewCollection_IsReadyToUse(t *testing.T) {
	factory := func(pkg *common.DataRecord, stream OutStream) (CSVWriter, error) {
		buf := bytes.NewBuffer([]byte{})
		csv := NewCSV(buf, "test")
		return csv, nil
	}
	got := NewCollection(factory)
	err := got.Write(
		&common.DataRecord{Data: &aez.CCDImage{PackData: &aez.CCDImagePackData{}}},
	)
	if err != nil {
		t.Errorf("NewCollection() returned collection that couldn't write, %v", err)
	}
}

func TestCSVCollection_Write(t *testing.T) {
	tests := []struct {
		name        string
		pkgs        []common.DataRecord
		wantErr     bool
		wantLines   int
		openStreams []OutStream
	}{
		{
			"Unknown stream should not write",
			[]common.DataRecord{{}},
			false,
			0,
			[]OutStream{},
		},
		{
			"CCDImage -> CCD",
			[]common.DataRecord{{Data: &aez.CCDImage{PackData: &aez.CCDImagePackData{}}}},
			false,
			3,
			[]OutStream{CCD},
		},
		{
			"Writing twice adds just one more line",
			[]common.DataRecord{
				{Data: &aez.CCDImage{PackData: &aez.CCDImagePackData{}}},
				{Data: &aez.CCDImage{PackData: &aez.CCDImagePackData{}}},
			},
			false,
			4,
			[]OutStream{CCD},
		},
		{
			"Two different streams",
			[]common.DataRecord{{Data: &aez.HTR{}}, {Data: &aez.STAT{}}},
			false,
			6, //Our simple factory puts everything in same buffer
			[]OutStream{HTR, STAT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})
			factory := func(pkg *common.DataRecord, stream OutStream) (CSVWriter, error) {
				csv := NewCSV(buf, "test")
				return csv, nil
			}
			col := NewCollection(factory)
			for _, pkg := range tt.pkgs {
				err := col.Write(&pkg)
				if err != nil {
					t.Errorf("CSVCollection.Write() = %v, wantErr %v", err, tt.wantErr)
				}
			}
			openStreams := make([]string, 0, len(col.streams))
			for k := range col.streams {
				openStreams = append(openStreams, k.String())
			}
			for _, stream := range tt.openStreams {
				_, ok := col.streams[stream]
				if !ok {
					t.Errorf(
						"CSVCollection.Write() caused %v streams, expected %v open",
						openStreams,
						stream.String(),
					)
				}
			}
			if len(openStreams) != len(tt.openStreams) {
				t.Errorf("CSVCollection has %v streams open, want %v", openStreams, tt.openStreams)
			}
			col.CloseAll()
			if lines := strings.Count(string(buf.Bytes()), "\n"); lines != tt.wantLines {
				t.Errorf(
					"CSVCollection.Write() added %v lines to output, want %v",
					lines,
					tt.wantLines,
				)
			}
		})
	}
}
