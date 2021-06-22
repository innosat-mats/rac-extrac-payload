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
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func Test_makeUnfinishedMultiPackError(t *testing.T) {
	// NON-STANDARD test since comparison of errors not recommended
	type args struct {
		multiPackBuffer *bytes.Buffer
		sourcePacket    common.DataRecord
	}
	tests := []struct {
		name       string
		args       args
		wantErr    string
		wantBuffer []byte
	}{
		{
			"Adds error and buffer to DataRecord",
			args{
				multiPackBuffer: bytes.NewBuffer([]byte("Hello")),
				sourcePacket: common.DataRecord{
					Origin: &common.OriginDescription{Name: "myname"},
					RID:    aez.CCD4,
				},
			},
			"orphaned multi-package data without termination detected [myname]",
			[]byte("Hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeUnfinishedMultiPackError(tt.args.multiPackBuffer, tt.args.sourcePacket)
			if got.Error.Error() != tt.wantErr {
				t.Errorf("makeUnfinishedMultiPackError().Error = %v, want %v", got.Error, tt.wantErr)
			}
			if fmt.Sprintf("%v", got.Buffer) != fmt.Sprintf("%v", tt.wantBuffer) {
				t.Errorf(
					"makeUnfinishedMultiPackError().Buffer = %v, want %v",
					got.Buffer,
					tt.wantBuffer,
				)
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
			"Three continuation packages makes one error for lacking end",
			[]common.DataRecord{
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42"),
				},
			},
			[]outcome{
				// For the unexpected start error we don't care about removing SID/RID
				{wantErr: true, partialErrMsg: "dangling final multipacket"},
			},
		},
		{
			"Returns standalone as it is",
			[]common.DataRecord{{
				Origin:       &common.OriginDescription{},
				SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0xc000},
				Buffer:       []byte("Hello"),
			}},
			[]outcome{{wantErr: false, bufferLength: 5}},
		},
		{
			"Returns aggregated",
			[]common.DataRecord{
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42 "),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x0000},
					Buffer:       []byte("42World"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x8000},
					Buffer:       []byte("42!"),
				},
			},
			[]outcome{{wantErr: false, bufferLength: 12}},
		},
		{
			"Errors if already started then continues multi",
			[]common.DataRecord{
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x8000},
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
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0x4000},
					Buffer:       []byte("Hello"),
				},
				{
					Origin:       &common.OriginDescription{},
					SourceHeader: &innosat.SourcePacketHeader{PacketSequenceControl: 0xc000},
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

func Test_makePackageInfo(t *testing.T) {
	tests := []struct {
		name         string
		sourcePacket *common.DataRecord
		want         string
	}{
		{
			"Empty",
			&common.DataRecord{},
			"[]",
		},
		{
			"With name",
			&common.DataRecord{
				Origin: &common.OriginDescription{
					Name: "test.rac",
				},
			},
			"[test.rac]",
		},
		{
			"With name & SourceHeader",
			&common.DataRecord{
				Origin: &common.OriginDescription{
					Name: "test.rac",
				},
				SourceHeader: &innosat.SourcePacketHeader{
					PacketID: 42,
				},
			},
			"[test.rac / Packet ID 42]",
		},
		{
			"With name & SourceHeader & RamsesTHHeader",
			&common.DataRecord{
				Origin: &common.OriginDescription{
					Name: "test.rac",
				},
				SourceHeader: &innosat.SourcePacketHeader{
					PacketID: 42,
				},
				RamsesTMHeader: &ramses.TMHeader{
					VCFrameCounter: 12,
				},
			},
			"[test.rac / Packet ID 42 / VC Frame Counter 12]",
		},
		{
			"With name & SourceHeader & RamsesTHHeader & RamsesHeader",
			&common.DataRecord{
				Origin: &common.OriginDescription{
					Name: "test.rac",
				},
				SourceHeader: &innosat.SourcePacketHeader{
					PacketID: 42,
				},
				RamsesTMHeader: &ramses.TMHeader{
					VCFrameCounter: 12,
				},
				RamsesHeader: &ramses.Ramses{
					Date: 666,
					Time: 420,
				},
			},
			"[test.rac / Packet ID 42 / VC Frame Counter 12 / Date 666, Time 420]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makePackageInfo(tt.sourcePacket); got != tt.want {
				t.Errorf("makePackageInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
