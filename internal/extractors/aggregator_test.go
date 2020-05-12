package extractors

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func Test_makeUnfinishedMultiPackError(t *testing.T) {
	// NON-STANDARD test since comparison of errors not recommended
	type args struct {
		multiPackBuffer *bytes.Buffer
		sourcePacket    common.DataRecord
	}
	tests := []struct {
		name string
		args args
		want common.DataRecord
	}{
		{
			"Adds error and buffer to DataRecord",
			args{
				multiPackBuffer: bytes.NewBuffer([]byte("Hello")),
				sourcePacket:    common.DataRecord{RID: aez.CCD4},
			},
			common.DataRecord{
				RID:    aez.CCD4,
				Error:  fmt.Errorf("orphaned multi-package data without termination detected"),
				Buffer: []byte("Hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeUnfinishedMultiPackError(tt.args.multiPackBuffer, tt.args.sourcePacket)
			if got.Error.Error() != tt.want.Error.Error() {
				t.Errorf("makeUnfinishedMultiPackError().Error = %v, want %v", got.Error, tt.want.Error)
			}
			if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", tt.want) {
				t.Errorf("makeUnfinishedMultiPackError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator(t *testing.T) {
	// NON-STANDARD test required to check control flow
	type outcome struct {
		wantErr       bool
		partialErrMsg string
		bufferLength  int
	}
	tests := []struct {
		name          string
		sourcePackets []common.DataRecord
		outcomes      []outcome
	}{
		{
			"Propagates errors without caring for content",
			[]common.DataRecord{{Error: errors.New("Hello")}},
			[]outcome{{wantErr: true, partialErrMsg: "Hello"}},
		},
		{
			"Three continuation packages makes one error for lacking start and one for lacking end",
			[]common.DataRecord{{Buffer: []byte("42")}, {Buffer: []byte("42")}, {Buffer: []byte("42")}},
			[]outcome{
				// For the unexpected start error we don't care about removing SID/RID
				{wantErr: true, partialErrMsg: "continuation packet", bufferLength: 2},
				{wantErr: true, partialErrMsg: "dangling final multipacket"},
			},
		},
		{
			"Returns standalone as it is",
			[]common.DataRecord{{
				SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0xc000},
				Buffer:       []byte("Hello"),
			}},
			[]outcome{{wantErr: false, bufferLength: 5}},
		},
		{
			"Returns aggregated",
			[]common.DataRecord{
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42 "),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42World"),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x8000},
					Buffer:       []byte("42!"),
				},
			},
			[]outcome{{wantErr: false, bufferLength: 12}},
		},
		{
			"Errors if already started then continues multi",
			[]common.DataRecord{
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x8000},
					Buffer:       []byte("42World!"),
				},
			},
			[]outcome{
				{wantErr: true, partialErrMsg: "orphaned", bufferLength: 5},
				{wantErr: false, bufferLength: 11},
			},
		},
		{
			"Errors if already started then reports standalone",
			[]common.DataRecord{
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					SourceHeader: innosat.SourcePacketHeader{PacketSequenceControl: 0xc000},
					Buffer:       []byte("World!"),
				},
			},
			[]outcome{
				{wantErr: true, partialErrMsg: "orphaned", bufferLength: 5},
				{wantErr: false, bufferLength: 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			sourceChan := make(chan common.DataRecord)
			targetChan := make(chan common.DataRecord)

			// Run aggregator
			go Aggregator(targetChan, sourceChan)

			// Fill up queue
			go func(source chan<- common.DataRecord) {
				defer close(source)
				for _, sourcePacket := range tt.sourcePackets {
					source <- sourcePacket
				}
			}(sourceChan)

			wg.Add(1)

			// Test resulting DataRecords
			go func(target <-chan common.DataRecord) {
				var nthTarget int
				for got := range target {
					if nthTarget >= len(tt.outcomes) {
						t.Errorf(
							"Outcome %v: Got unexpected DataRecord %v, only wanted %v records",
							nthTarget,
							got,
							len(tt.outcomes),
						)
					} else {
						outcome := tt.outcomes[nthTarget]
						if outcome.wantErr && got.Error == nil {
							t.Errorf("Outcome %v: Wanted error '%v', found nothing", nthTarget, outcome.partialErrMsg)
						} else if !outcome.wantErr && got.Error != nil {
							t.Errorf("Outcome %v: Unexpected error '%v', wanted nothing", nthTarget, got.Error)
						} else if outcome.wantErr && !strings.Contains(got.Error.Error(), outcome.partialErrMsg) {
							t.Errorf("Outcome %v: Error '%v' does not contain '%v'", nthTarget, got.Error.Error(), outcome.partialErrMsg)
						}
						if len(got.Buffer) != outcome.bufferLength {
							t.Errorf(
								"Outcome %v: Expected buffer size %v, found %v",
								nthTarget,
								outcome.bufferLength,
								len(got.Buffer),
							)
						}
					}

					nthTarget++
				}
				if nthTarget != len(tt.outcomes) {
					t.Errorf(
						"Expected %v DataRecords produced, but found %v",
						len(tt.outcomes),
						nthTarget,
					)
				}
				wg.Done()
			}(targetChan)

			wg.Wait()
		})
	}
}
