package extractors

import (
	"os"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func TestDregs_getDregsFileName(t *testing.T) {
	type fields struct {
		Path    string
		MaxDiff int64
	}
	type args struct {
		data common.DataRecord
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"Creates correct filename",
			fields{"/some/path", MaxDeviationNanos},
			args{common.DataRecord{
				TMHeader: &innosat.TMHeader{
					CUCTimeSeconds:  42,
					CUCTimeFraction: 0,
				},
			}},
			"/some/path/42000000000.dregs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dregs := &Dregs{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			if got := dregs.getDregsFileName(tt.args.data); string(got) != string(tt.want) {
				t.Errorf("Dregs.getDregsFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDregs_DumpDregs(t *testing.T) {
	type fields struct {
		Path    string
		MaxDiff int64
	}
	type args struct {
		data common.DataRecord
	}
	dir, err := os.MkdirTemp("", "innosat-mats-dregs-dir")
	if err != nil {
		t.Errorf("TestDregs_DumpDregs() could not setup output directory '%v'", err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Creates expected file",
			fields{dir, MaxDeviationNanos},
			args{common.DataRecord{
				TMHeader: &innosat.TMHeader{
					CUCTimeSeconds:  42,
					CUCTimeFraction: 0,
				},
			}},
			false,
		},
		{
			"Returns error if no path set",
			fields{},
			args{common.DataRecord{}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dregs := &Dregs{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			err := dregs.DumpDregs(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dregs.DumpDregs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				dregsFile := dregs.getDregsFileName(tt.args.data)
				_, err = os.ReadFile(dregsFile)
				if err != nil {
					t.Errorf(
						"Dregs.DumpDregs() expected to produce file '%v', but got error reading it: %v",
						dregsFile,
						err,
					)
				}
			}
			if tt.wantErr && err != ErrNoDregsPath {
				t.Errorf(
					"Dregs.DumpDregs() was expected to return %v but got %v",
					ErrNoDregsPath,
					err,
				)
			}
		})
	}
}

func TestDregs_GetDregs(t *testing.T) {
	type fields struct {
		Path    string
		MaxDiff int64
	}
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"Reads the right file",
			fields{
				Path:    "./dregs",
				MaxDiff: 10,
			},
			args{50},
			[]byte("This is the right file."),
			false,
		},
		{
			"Returns error if no path set",
			fields{},
			args{50},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dregs := &Dregs{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			got, err := dregs.GetDregs(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dregs.GetDregs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dregs.GetDregs() = '%v', want '%v'", string(got), string(tt.want))
			}
			if tt.wantErr && err != ErrNoDregsPath {
				t.Errorf(
					"Dregs.GetDregs() was expected to return %v but got %v",
					ErrNoDregsPath,
					err,
				)
			}
		})
	}
}
