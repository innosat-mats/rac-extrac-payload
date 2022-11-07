package extractors

import (
	"os"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func TestSlask_getSlaskFileName(t *testing.T) {
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
			"/some/path/42000000000.slask",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slask := &Slask{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			if got := slask.getSlaskFileName(tt.args.data); string(got) != string(tt.want) {
				t.Errorf("Slask.getSlaskFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlask_DumpSlask(t *testing.T) {
	type fields struct {
		Path    string
		MaxDiff int64
	}
	type args struct {
		data common.DataRecord
	}
	dir, err := os.MkdirTemp("", "innosat-mats-slask-dir")
	if err != nil {
		t.Errorf("TestSlask_DumpSlask() could not setup output directory '%v'", err)
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
			slask := &Slask{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			err := slask.DumpSlask(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slask.DumpSlask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				slaskFile := slask.getSlaskFileName(tt.args.data)
				_, err = os.ReadFile(slaskFile)
				if err != nil {
					t.Errorf(
						"Slask.DumpSlask() expected to produce file '%v', but got error reading it: %v",
						slaskFile,
						err,
					)
				}
			}
			if tt.wantErr && err != ErrNoSlaskPath {
				t.Errorf(
					"Slask.DumpSlask() was expected to return %v but got %v",
					ErrNoSlaskPath,
					err,
				)
			}
		})
	}
}

func TestSlask_GetSlask(t *testing.T) {
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
			"Returns error if no path set",
			fields{
				Path:    "./slask",
				MaxDiff: 10,
			},
			args{50},
			[]byte("This is the right file.\n"),
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
			slask := &Slask{
				Path:    tt.fields.Path,
				MaxDiff: tt.fields.MaxDiff,
			}
			got, err := slask.GetSlask(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slask.GetSlask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slask.GetSlask() = '%v', want '%v'", string(got), string(tt.want))
			}
			if tt.wantErr && err != ErrNoSlaskPath {
				t.Errorf(
					"Slask.GetSlask() was expected to return %v but got %v",
					ErrNoSlaskPath,
					err,
				)
			}
		})
	}
}
