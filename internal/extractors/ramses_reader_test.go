package extractors

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

var ramesesHeaderSize int = binary.Size(ramses.Ramses{})
var tmHeaderSize int = binary.Size(ramses.TMHeader{})
var minTotalSize int = ramesesHeaderSize + tmHeaderSize

func TestGetRecord(t *testing.T) {
	type args struct {
		buf io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"all zeros, not a vaild file",
			args{bytes.NewReader(make([]byte, minTotalSize))},
			[]byte{},
			true,
		},
		{
			"too short buffer",
			args{bytes.NewReader([]byte{0x90, 0xeb, 0})},
			[]byte{},
			true,
		},
		{
			"a correct package 0 length no payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 16, 0}, make([]byte, minTotalSize-4)...),
				),
			},
			[]byte{},
			false,
		},
		{
			"a correct package 1 byte payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 17, 0}, make([]byte, minTotalSize-4+1)...),
				),
			},
			[]byte{0},
			false,
		},
		{
			"a package with too short payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 18, 0}, make([]byte, minTotalSize-4+1)...),
				),
			},
			[]byte{},
			true,
		},
		{
			"a package without full OHBSE CCSDS TM Packet Header",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 16, 0}, make([]byte, minTotalSize-5)...),
				),
			},
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := getRecord(StreamBatch{
				Origin: &common.OriginDescription{Name: "myname.rac"},
				Buf:    tt.args.buf,
			})
			if !reflect.DeepEqual(got.Buffer, tt.want) {
				t.Errorf("Package.Payload = %v, want %v", got.Buffer, tt.want)
			}
			if (got.Error != nil) != tt.wantErr {
				t.Errorf("getRecord() error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}

		})
	}
}

func TestGetRecord_FromFile(t *testing.T) {
	type params struct {
		fileContent []byte
		wantDone    bool
		wantErr     bool
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			"Too short file",
			params{
				fileContent: append([]byte{0x90, 0xeb, 18, 0}, make([]byte, minTotalSize-4+1)...),
				wantDone:    false,
				wantErr:     true,
			},
		},
		{
			"No more record / empty file",
			params{
				fileContent: []byte{},
				wantDone:    true,
				wantErr:     false,
			},
		},
	}
	for _, tt := range tests {
		file, err := os.CreateTemp("", "test-file")
		defer os.Remove(file.Name())
		if err != nil {
			t.Errorf("Could not create testfile %v", err)
			return
		}
		file.Write(tt.params.fileContent)
		file.Seek(0, 0)
		got, done := getRecord(StreamBatch{
			Origin: &common.OriginDescription{Name: "myname.rac"},
			Buf:    file,
		})
		if done != tt.params.wantDone {
			t.Errorf("getRecord() done = %v, want %v", done, tt.params.wantDone)
		}
		if (got.Error != nil) != tt.params.wantErr {
			t.Errorf("getRecord() error = %v, wantError %v", got.Error, tt.params.wantErr)
		}
	}
}

func TestDecodeRamses(t *testing.T) {
	type streams struct {
		buf    []byte
		origin common.OriginDescription
	}
	type outcome struct {
		wantErr    bool
		originName string
		buf        []byte
	}
	tests := []struct {
		name     string
		streams  []streams
		outcomes []outcome
	}{
		{
			"Works with no input",
			[]streams{},
			[]outcome{},
		},
		{
			"Emtpy stream gives no output",
			[]streams{{}},
			[]outcome{},
		},
		{
			"Terminates on first invalid record",
			[]streams{
				{
					origin: common.OriginDescription{Name: "No. 1"},
					buf: append(append(
						//OK record
						append([]byte{0x90, 0xeb, 17, 0}, make([]byte, minTotalSize-4+1)...),
						//Invalid record but correct size
						append([]byte{0x42, 0xeb, 17, 0}, make([]byte, minTotalSize-4+1)...)...,
					),
						//OK record
						append([]byte{0x90, 0xeb, 17, 0}, make([]byte, minTotalSize-4+1)...)...,
					),
				},
			},
			[]outcome{
				{
					originName: "No. 1",
					buf:        []byte{0},
				},
				{
					originName: "No. 1",
					buf:        []byte{},
					wantErr:    true,
				},
			},
		},
		{
			"Reports errors first",
			[]streams{
				{
					// OK Stream
					origin: common.OriginDescription{Name: "No. 1"},
					buf:    append([]byte{0x90, 0xeb, 17, 0}, make([]byte, minTotalSize-4+1)...),
				},
				{
					// BAD Stream
					origin: common.OriginDescription{Name: "No. 2"},
					buf:    []byte{0x90, 0xeb, 17, 0},
				},
			},
			[]outcome{
				{
					wantErr:    true,
					originName: "No. 2",
				},
				{
					originName: "No. 1",
					buf:        []byte{0},
				},
			},
		},
		{
			"Sorts on first packet per stream",
			[]streams{
				{
					// OK Stream, 2 Packets (10 days, 42 millis)
					origin: common.OriginDescription{Name: "No. 1"},
					buf: append(
						append(
							[]byte{0x90, 0xeb, 17, 0, 0, 0, 0, 0, 42, 0, 0, 0, 10, 0, 0, 0},
							make([]byte, tmHeaderSize+1)...,
						),
						append(
							[]byte{0x90, 0xeb, 18, 0, 0, 0, 0, 0, 42, 0, 0, 0, 10, 0, 0, 0},
							make([]byte, tmHeaderSize+2)...,
						)...,
					),
				},
				{
					// OK Stream, 2 Packets earlier date (10 days 41 millis)
					origin: common.OriginDescription{Name: "No. 2"},
					buf: append(
						append(
							[]byte{0x90, 0xeb, 17, 0, 0, 0, 0, 0, 41, 0, 0, 0, 10, 0, 0, 0},
							make([]byte, tmHeaderSize+1)...,
						),
						append(
							[]byte{0x90, 0xeb, 18, 0, 0, 0, 0, 0, 41, 0, 0, 0, 10, 0, 0, 0},
							make([]byte, tmHeaderSize+2)...,
						)...,
					),
				},
			},
			[]outcome{
				{
					originName: "No. 2",
					buf:        []byte{0},
				},
				{
					originName: "No. 2",
					buf:        []byte{0, 0},
				},
				{
					originName: "No. 1",
					buf:        []byte{0},
				},
				{
					originName: "No. 1",
					buf:        []byte{0, 0},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packs := make(chan common.DataRecord)
			streams := make([]StreamBatch, len(tt.streams))
			for i := range streams {
				streams[i] = StreamBatch{
					Buf:    bytes.NewReader(tt.streams[i].buf),
					Origin: &tt.streams[i].origin,
				}
			}
			go DecodeRamses(packs, streams...)
			var idxOutcome int = -1
			for got := range packs {
				idxOutcome++
				if idxOutcome >= len(tt.outcomes) {
					t.Errorf(
						"%v: Got unexpected outcome %v, only wanted %v",
						idxOutcome,
						got,
						len(tt.outcomes),
					)
					continue
				}
				if (got.Error != nil) != tt.outcomes[idxOutcome].wantErr {
					t.Errorf(
						"%v: DataRecord.Error = %v, wantErr %v",
						idxOutcome,
						got.Error,
						tt.outcomes[idxOutcome].wantErr,
					)
				}
				if got.Error == nil && !reflect.DeepEqual(got.Buffer, tt.outcomes[idxOutcome].buf) {
					t.Errorf(
						"%v DataRecord.Buffer = %v, want %v",
						idxOutcome,
						got.Buffer,
						tt.outcomes[idxOutcome].buf,
					)
				}
				if !reflect.DeepEqual(got.Origin.Name, tt.outcomes[idxOutcome].originName) {
					t.Errorf(
						"%v DataRecord.Origin.Name = %v, want %v",
						idxOutcome,
						got.Origin.Name,
						tt.outcomes[idxOutcome].originName,
					)
				}
			}
			if idxOutcome+1 != len(tt.outcomes) {
				t.Errorf("Got %v DataRecords, want %v", idxOutcome+1, len(tt.outcomes))
			}
		})
	}
}

func TestStreamBatch_OriginName(t *testing.T) {
	tests := []struct {
		name   string
		origin *common.OriginDescription
		want   string
	}{
		{"Default empty", nil, ""},
		{"Returns name", &common.OriginDescription{Name: "Test"}, "Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := &StreamBatch{Origin: tt.origin}
			if got := stream.OriginName(); got != tt.want {
				t.Errorf("StreamBatch.OriginName() = %v, want %v", got, tt.want)
			}
		})
	}
}
