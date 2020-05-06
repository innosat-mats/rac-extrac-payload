package exports

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_csvName(t *testing.T) {
	type args struct {
		dir        string
		originName string
		packetType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Case 1", args{".", "somefile.rac", "TEST"}, "somefile_TEST.csv"},
		{"Case 2", args{"my/dir", "somefile", "TEST"}, "my/dir/somefile_TEST.csv"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := csvName(tt.args.dir, tt.args.originName, tt.args.packetType); got != tt.want {
				t.Errorf("csvName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiskCallbackFactory(t *testing.T) {
	type args struct {
		writeImages     bool
		writeTimeseries bool
	}
	type wantFile struct {
		base  string
		lines int
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
				{"File1_STAT.csv", 4},
			},
		},
		{
			"Swaps file if different origin",
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
				{"File1_STAT.csv", 3},
				{"File2_STAT.csv", 3},
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
					Data:   aez.CPRU{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.CPRU{},
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
			},
			[]wantFile{
				{"File1_STAT.csv", 3},
				{"File1_CPRU.csv", 4},
				{"File1_HTR.csv", 4},
				{"File1_PWR.csv", 4},
			},
		},
		{
			"Swaps out the other files",
			args{writeTimeseries: true},
			[]common.DataRecord{
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.CPRU{},
				},
				{
					Origin: common.OriginDescription{Name: "File2.rac"},
					Data:   aez.CPRU{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.HTR{},
				},
				{
					Origin: common.OriginDescription{Name: "File2.rac"},
					Data:   aez.HTR{},
				},
				{
					Origin: common.OriginDescription{Name: "File1.rac"},
					Data:   aez.PWR{},
				},
				{
					Origin: common.OriginDescription{Name: "File2.rac"},
					Data:   aez.PWR{},
				},
			},
			[]wantFile{
				{"File1_CPRU.csv", 3},
				{"File1_HTR.csv", 3},
				{"File1_PWR.csv", 3},
				{"File2_CPRU.csv", 3},
				{"File2_HTR.csv", 3},
				{"File2_PWR.csv", 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup and cleanup of output directory
			dir, err := ioutil.TempDir("/tmp", "innosat-mats")
			if err != nil {
				t.Errorf("DiskCallbackFactory() could not setup output directory '%v'", err)
			}
			defer os.RemoveAll(dir)

			// Produce callback and teardown
			callback, teardown := DiskCallbackFactory(dir, tt.args.writeImages, tt.args.writeTimeseries)

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
				if newLines := strings.Count(string(content), "\n"); newLines != want.lines {
					t.Errorf("DiskCallbackFactory() expected file %v to have %v lines, found %v", want.base, want.lines, newLines)
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
