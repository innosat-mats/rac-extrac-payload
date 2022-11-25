package exports

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
	"github.com/innosat-mats/rac-extract-payload/internal/timeseries"
)

func Test_parquetName(t *testing.T) {
	type args struct {
		dir    string
		packet common.DataRecord
		stream timeseries.OutStream
	}
	record := common.DataRecord{
		Origin:         &common.OriginDescription{Name: "File1.rac"},
		RamsesHeader:   &ramses.Ramses{},
		RamsesTMHeader: &ramses.TMHeader{},
		SourceHeader:   &innosat.SourcePacketHeader{},
		TMHeader:       &innosat.TMHeader{},
		Data:           &aez.STAT{},
	}
	stream := timeseries.OutStreamFromDataRecord(&record)
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Case 1", args{".", record, stream}, filepath.FromSlash("1980/1/5/STAT_File1.parquet")},
		{"Case 2", args{"my/dir", record, stream}, filepath.FromSlash("my/dir/1980/1/5/STAT_File1.parquet")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parquetName(tt.args.dir, &tt.args.packet, tt.args.stream); got != tt.want {
				t.Errorf("parquetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParquetCallbackFactoryCreator(t *testing.T) {
	type args struct {
		wg *sync.WaitGroup
	}
	type wantFile struct {
		base string
	}
	tests := []struct {
		name         string
		args         args
		callbackArgs []common.DataRecord
		wantFiles    []wantFile
	}{

		{
			"Appends to open file if same origin",
			args{},
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
			},
			[]wantFile{
				{"STAT_File1.parquet"},
			},
		},
		{
			"Adds from all racs into same file",
			args{},
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File2.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
			},
			[]wantFile{
				{"STAT_File1.parquet"},
				{"STAT_File2.parquet"},
			},
		},
		{
			"Handles all types in parallel",
			args{},
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.STAT{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.CPRU{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.CPRU{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.HTR{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.HTR{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.PWR{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.PWR{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.PMData{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.PMData{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.TCAcceptSuccessData{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.TCExecSuccessData{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.TCAcceptFailureData{},
				},
				{
					Origin:         &common.OriginDescription{Name: "File1.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.TCExecFailureData{},
				},
			},
			[]wantFile{
				{"CPRU_File1.parquet"},
				{"HTR_File1.parquet"},
				{"PM_File1.parquet"},
				{"PWR_File1.parquet"},
				{"STAT_File1.parquet"},
				{"TCV_File1.parquet"},
			},
		},
		{
			"Continues on error due to wrong image shape",
			args{},
			[]common.DataRecord{
				{
					TMHeader: &innosat.TMHeader{},
					RID:      aez.CCD2,
					Origin:   &common.OriginDescription{Name: "File1.rac"},
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  42,
							NROW:  42,
							EXPTS: 5,
						},
						ImageFileName: "File1_wrong_shape_2.png",
					},
					Buffer: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
				{
					TMHeader: &innosat.TMHeader{},
					RID:      aez.CCD3,
					Origin:   &common.OriginDescription{Name: "File2.rac"},
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  1,
							NROW:  2,
							EXPTS: 6,
						},
						ImageFileName: "File1_6000000000_3.png",
					},
					Buffer: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			[]wantFile{
				{"CCD_File1.parquet"},
				{"CCD_File2.parquet"},
			},
		},
		/*
		 */
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg = &sync.WaitGroup{}
			// Setup and cleanup of output directory
			dir, err := os.MkdirTemp("", "innosat-mats")
			if err != nil {
				t.Errorf("ParquetCallbackFactory() could not setup output directory '%v'", err)
			}

			// Produce callback and teardown
			callback, teardown := ParquetCallbackFactory(dir, tt.args.wg)

			// Invoke callback and then teardown
			for _, pkg := range tt.callbackArgs {
				callback(pkg)
			}
			teardown()

			savePath := filepath.Join(dir, "1980", "1", "5")
			for _, want := range tt.wantFiles {
				// Test each output for file name and expected number of lines
				path := filepath.Join(savePath, want.base)
				_, err := os.ReadFile(path)
				if err != nil {
					t.Errorf(
						"ParquetCallbackFactory() expected to produce file '%v', but got error reading it: %v",
						path, err,
					)
					files, _ := os.ReadDir(savePath)
					fmt.Printf("Files in %v:\n", savePath)
					for _, file := range files {
						fmt.Printf("  %v %v\n", file.Name(), file.IsDir())
					}
				}
			}

			// Test that number of output files equals expected
			files, err := os.ReadDir(savePath)
			if err != nil {
				t.Errorf("ParquetCallbackFactory() could not read directory: %v", err)
			}
			if nFiles, expect := len(files), len(tt.wantFiles); nFiles != expect {
				t.Errorf(
					"ParquetCallbackFactory() created %v files, expected %v files",
					nFiles, expect,
				)
				fmt.Printf("Files in %v:\n", savePath)
				for _, file := range files {
					fmt.Printf("  %v %v\n", file.Name(), file.IsDir())
				}
			}

		})
	}
}
