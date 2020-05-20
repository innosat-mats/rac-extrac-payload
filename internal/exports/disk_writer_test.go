package exports

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_csvName(t *testing.T) {
	type args struct {
		dir        string
		packetType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Case 1", args{".", "TEST"}, "TEST.csv"},
		{"Case 2", args{"my/dir", "TEST"}, filepath.FromSlash("my/dir/TEST.csv")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := csvName(tt.args.dir, tt.args.packetType); got != tt.want {
				t.Errorf("csvName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiskCallbackFactoryCreator(t *testing.T) {
	type args struct {
		writeImages     bool
		writeTimeseries bool
		wg              *sync.WaitGroup
	}
	type wantFile struct {
		base           string
		lines          int
		dontCountLines bool
	}
	tests := []struct {
		name         string
		args         args
		callbackArgs []common.DataRecord
		wantFiles    []wantFile
	}{
		{
			"Doesn't create files if no writeTimeseries",
			args{writeTimeseries: false},
			[]common.DataRecord{
				{Data: aez.STAT{}},
			},
			[]wantFile{},
		},
		{
			"Appends to open file if same origin",
			args{writeTimeseries: true},
			[]common.DataRecord{
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.STAT{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.STAT{},
				},
			},
			[]wantFile{
				{"STAT.csv", 4, false},
			},
		},
		{
			"Adds from all racs into same file",
			args{writeTimeseries: true},
			[]common.DataRecord{
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.STAT{},
				},
				{
					Origin: common.OriginDescription{Name: "File2.rac"},
					Data:   aez.STAT{},
				},
			},
			[]wantFile{
				{"STAT.csv", 4, false},
			},
		},
		{
			"Handles all types in parallel",
			args{writeTimeseries: true},
			[]common.DataRecord{
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.STAT{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   &aez.CPRU{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   &aez.CPRU{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.HTR{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.HTR{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.PWR{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.PWR{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.PMData{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.PMData{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.TCAcceptSuccessData{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.TCExecSuccessData{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.TCAcceptFailureData{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.TCExecFailureData{},
				},
			},
			[]wantFile{
				{"STAT.csv", 3, false},
				{"CPRU.csv", 4, false},
				{"HTR.csv", 4, false},
				{"PWR.csv", 4, false},
				{"PM.csv", 4, false},
				{"TCV.csv", 6, false},
			},
		},
		{
			"Creates images and jsons",
			args{writeImages: true},
			[]common.DataRecord{
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  1,
							NROW:  2,
							EXPTS: 5,
						},
					},
					Buffer: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  1,
							NROW:  2,
							EXPTS: 6,
						},
					},
					Buffer: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			[]wantFile{
				{"File1_5000000000.png", 0, true},
				{"File1_5000000000.json", 0, true},
				{"File1_6000000000.png", 0, true},
				{"File1_6000000000.json", 0, true},
			},
		},
		{
			"Doesn't creates images when asked not to",
			args{writeImages: false},
			[]common.DataRecord{
				{
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  1,
							NROW:  2,
							EXPTS: 5,
						},
					},
					Buffer: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			[]wantFile{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg = &sync.WaitGroup{}
			// Setup and cleanup of output directory
			dir, err := ioutil.TempDir("", "innosat-mats")
			if err != nil {
				t.Errorf("DiskCallbackFactory() could not setup output directory '%v'", err)
			}
			defer os.RemoveAll(dir)

			// Produce callback and teardown
			callback, teardown := DiskCallbackFactory(dir, tt.args.writeImages, tt.args.writeTimeseries, tt.args.wg)

			// Invoke callback and then teardown
			for _, pkg := range tt.callbackArgs {
				callback(pkg)
			}
			teardown()

			for _, want := range tt.wantFiles {
				// Test each output for file name and expected number of lines
				path := filepath.Join(dir, want.base)
				content, err := ioutil.ReadFile(path)
				if err != nil {
					t.Errorf("DiskCallbackFactory() expected to produce file '%v', but got error reading it: %v", path, err)
				}
				if !want.dontCountLines {
					if newLines := strings.Count(string(content), "\n"); newLines != want.lines {
						t.Errorf("DiskCallbackFactory() expected file %v to have %v lines, found %v", want.base, want.lines, newLines)
					}
				}
			}

			// Test that number of output files equals expected
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				t.Errorf("DiskCallbackFactory() could not read directory: %v", err)
			}
			if nFiles, expect := len(files), len(tt.wantFiles); nFiles != expect {
				t.Errorf("DiskCallbackFactory() created %v files, expected %v files", nFiles, expect)
			}

		})
	}
}
