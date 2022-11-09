package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/extractors"
)

func Test_getCallback(t *testing.T) {
	type args struct {
		toStdout       bool
		toParquet      bool
		toAws          bool
		project        string
		skipImages     bool
		skipTimeseries bool
		wg             *sync.WaitGroup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Returns stdout callback", args{toStdout: true}, false},
		{"Returns aws callback", args{toAws: true, project: "test"}, false},
		{"Returns aws callback requires project", args{toAws: true}, true},
		{"Returns disk callback", args{project: "somewhere"}, false},
		{"Returns error if no output directory", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := getCallback(
				tt.args.toStdout,
				tt.args.toParquet,
				tt.args.toAws,
				tt.args.project,
				tt.args.skipImages,
				tt.args.skipTimeseries,
				"",
				tt.args.wg,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_processFiles(t *testing.T) {
	type args struct {
		inputFiles []string
		callback   common.Callback
	}
	type fixtures struct {
		files []string
	}
	tests := []struct {
		name     string
		args     args
		fixtures fixtures
		wantErr  bool
	}{
		{
			"Processes all files",
			args{inputFiles: []string{"a", "b"}},
			fixtures{files: []string{"a", "b", "c"}},
			false,
		},
		{
			"Fails on missing files",
			args{inputFiles: []string{"a", "b"}},
			fixtures{files: []string{}},
			true,
		},
	}
	mapFilenamesToDirectory := func(dir string, files []string) []string {
		newFiles := make([]string, len(files))
		for idx, file := range files {
			newFiles[idx] = filepath.Join(dir, file)
		}
		return newFiles
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "mats-testing")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)
			for _, fileName := range tt.fixtures.files {
				filePath := filepath.Join(dir, fileName)
				file, err := os.Create(filePath)
				if err != nil {
					log.Fatal(err)
				}
				file.Close()
			}
			updatedFilenames := mapFilenamesToDirectory(dir, tt.args.inputFiles)
			extractor := func(
				callback common.Callback,
				dregs extractors.Dregs,
				streamBatch ...extractors.StreamBatch,
			) {
				ptCallback := reflect.ValueOf(callback).Pointer()
				ptArgsCallback := reflect.ValueOf(tt.args.callback).Pointer()
				if ptCallback != ptArgsCallback {
					t.Errorf("Expected callback to be passed on to extractor, got %v", callback)
				}

				if len(streamBatch) != len(tt.args.inputFiles) {
					t.Errorf("Expected %v streams but got %v", len(tt.args.inputFiles), len(streamBatch))
				}
				for idx, stream := range streamBatch {
					if stream.Buf == nil {
						t.Errorf("Expected stream %v to have a buffer but got nil", idx)
					}
					if stream.Origin.Name != updatedFilenames[idx] {
						t.Errorf(
							"Expected stream %v to have Name %v but got %v",
							idx,
							updatedFilenames[idx],
							stream.Origin.Name,
						)
					}
					elapsed := time.Since(stream.Origin.ProcessingDate)
					if elapsed < 0 || elapsed > time.Second {
						t.Errorf("Processing time %v seems not to be now", stream.Origin.ProcessingDate)
					}
				}
			}

			err = processFiles(
				extractor,
				updatedFilenames,
				extractors.Dregs{},
				tt.args.callback,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
