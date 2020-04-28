package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_generateStdoutCallback(t *testing.T) {
	type args struct {
		writeTimeseries bool
	}
	type innerArgs struct {
		dataRecord common.DataRecord
	}
	tests := []struct {
		name      string
		args      args
		innerArgs innerArgs
		want      string
	}{
		{"Prints nothing when not asked to", args{false}, innerArgs{common.DataRecord{}}, ""},
		{
			"Prints nothing when not asked to",
			args{true},
			innerArgs{common.DataRecord{}},
			"{Origin:{Name: ProcessingDate:0001-01-01 00:00:00 +0000 UTC} RamsesHeader:{Synch:0 Length:0 Port:0 Type:0 Secure:0 Time:0 Date:0} RamsesSecure:{IPAddress:0 Port:0 Seq:0 Retransmission:0 Ack:0 _:0} SourceHeader:{PacketID:0 PacketSequenceControl:0 PacketLength:0} TMHeader:{PUS:0 ServiceType:0 ServiceSubType:0 CUCTimeSeconds:0 CUCTimeFraction:0} Data:<nil> Error:<nil> Buffer:[]}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			callback := generateStdoutCallback(buf, tt.args.writeTimeseries)
			callback(tt.innerArgs.dataRecord)
			if got := buf.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateStdoutCallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCallback(t *testing.T) {
	type args struct {
		stdout          bool
		outputDirectory string
		skipImages      bool
		skipTimeseries  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Returns stdout callback", args{stdout: true}, false},
		{"Returns disk callback", args{outputDirectory: "somewhere"}, false},
		{"Returns error if no output directory", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getCallback(tt.args.stdout, tt.args.outputDirectory, tt.args.skipImages, tt.args.skipTimeseries)
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
		callback   common.ExtractCallback
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
			dir, err := ioutil.TempDir("/tmp", "mats-testing")
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
			extractor := func(callback common.ExtractCallback, streamBatch ...common.StreamBatch) {
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
					elapsed := time.Now().Sub(stream.Origin.ProcessingDate)
					if elapsed < 0 || elapsed > time.Second {
						t.Errorf("Processing time %v seems not to be now", stream.Origin.ProcessingDate)
					}
				}
			}

			err = processFiles(
				extractor,
				updatedFilenames,
				tt.args.callback,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
