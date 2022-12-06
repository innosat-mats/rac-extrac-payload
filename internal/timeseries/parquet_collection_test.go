package timeseries

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func TestParquetCollection_Write(t *testing.T) {
	tests := []struct {
		name        string
		pkgs        []common.DataRecord
		wantErr     bool
		openStreams []string
	}{
		{
			"Unknown stream should not write",
			[]common.DataRecord{{}},
			false,
			[]string{},
		},
		{
			"A single stream from same file",
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "test1"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
				{
					Origin:         &common.OriginDescription{Name: "test1"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
			},
			false,
			[]string{
				filepath.FromSlash("STAT/1980/1/5/test1.parquet"),
			},
		},
		{
			"Two different streams from same file",
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "test1"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.HTR{},
				},
				{
					Origin:         &common.OriginDescription{Name: "test1"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
			},
			false,
			[]string{
				filepath.FromSlash("HTR/1980/1/5/test1.parquet"),
				filepath.FromSlash("STAT/1980/1/5/test1.parquet"),
			},
		},
		{
			"Two different streams from different files",
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "test1"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
				{
					Origin:         &common.OriginDescription{Name: "test2"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
			},
			false,
			[]string{
				filepath.FromSlash("STAT/1980/1/5/test1.parquet"),
				filepath.FromSlash("STAT/1980/1/5/test2.parquet"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "innosat-mats")
			if err != nil {
				t.Errorf("ParquetCollection() could not setup output directory '%v'", err)
			}
			factory := func(pkg *common.DataRecord, stream OutStream) (ParquetWriter, error) {
				f := filepath.Join(dir, "test")
				parquet := NewParquet(f, pkg)
				return parquet, nil
			}
			col := NewParquetCollection(factory)
			for _, pkg := range tt.pkgs {
				err := col.Write(&pkg)
				if err != nil {
					t.Errorf("ParquetCollection.Write() = %v, wantErr %v", err, tt.wantErr)
				}
			}
			openStreams := make([]string, 0, len(col.streams))
			for k := range col.streams {
				openStreams = append(openStreams, k)
			}
			for _, stream := range tt.openStreams {
				_, ok := col.streams[stream]
				if !ok {
					t.Errorf(
						"ParquetCollection.Write() caused %v streams, expected %v open",
						openStreams,
						stream,
					)
				}
			}
			if len(openStreams) != len(tt.openStreams) {
				t.Errorf("ParquetCollection has %v streams open, want %v", openStreams, tt.openStreams)
			}
			col.CloseAll()
		})
	}
}
