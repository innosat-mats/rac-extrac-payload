package extractors

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

// makeInstrumentData is a test fixture creator for generating a byte slice of RID/SID the struct and any trailing bytes in the packet.
func makeInstrumentData(sidrid uint16, data interface{}, trailingBytes []byte) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, sidrid)
	binary.Write(&buf, binary.BigEndian, data)
	return append(buf.Bytes(), trailingBytes...)
}

func TestDecodeAEZ(t *testing.T) {

	tests := []struct {
		name       string
		arg        common.DataRecord
		want       common.Exportable
		wantSID    aez.SID
		wantRID    aez.RID
		wantBufLen int
		wantErr    bool
	}{
		{
			"Package with error",
			common.DataRecord{Error: io.EOF, Buffer: []byte("Hello")},
			nil,
			aez.SID(0),
			aez.RID(0),
			5,
			true,
		},

		{
			"STAT package",
			common.DataRecord{
				SourceHeader: innosat.SourcePacketHeader{PacketID: 0x0864, PacketSequenceControl: 0xc89a, PacketLength: 0x31},
				TMHeader:     innosat.TMDataFieldHeader{PUS: 16, ServiceType: 3, ServiceSubType: 0x19, CUCTimeSeconds: 0, CUCTimeFraction: 0},
				Buffer: []byte{0x00, 0x01, 0x7f, 0x04, 0x02, 0x82,
					0x04, 0x02, 0x02, 0x06, 0x01, 0x1b, 0x12, 0x00, 0x00, 0xef, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0xa5, 0xd5,
				},
			},
			aez.STAT{SPID: 32516, SPREV: 2, FPID: 33284, FPREV: 2, SVNA: 2, SVNB: 6, SVNC: 1, TS: 454164480, TSS: 61344, MODE: 2, EDACE: 0, EDACCE: 0, EDACN: 16777216, SPWEOP: 1090519040, SPWEEP: 0, ANOMALY: 0},
			aez.SIDSTAT,
			aez.RID(0),
			0,
			false,
		},
		{
			"Bad package",
			common.DataRecord{
				SourceHeader: innosat.SourcePacketHeader{PacketID: 0x0864, PacketSequenceControl: 0xc89a, PacketLength: 0x31},
				TMHeader:     innosat.TMDataFieldHeader{PUS: 16, ServiceType: 3, ServiceSubType: 0x19, CUCTimeSeconds: 0, CUCTimeFraction: 0},
				Buffer: []byte{0x00, 0x01, 0x7f, 0x04, 0x02, 0x82,
					0x04, 0x02, 0x02, 0x06, 0x01, 0x1b, 0x12, 0x00, 0x00, 0xef, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			aez.STAT{},
			aez.SIDSTAT,
			aez.RID(0),
			0,
			true,
		},
		{
			"Transparnt Data package",
			common.DataRecord{
				TMHeader: innosat.TMDataFieldHeader{ServiceType: 128, ServiceSubType: 25},
				Buffer: makeInstrumentData(
					uint16(aez.CCD2),
					aez.CCDImagePackData{NBC: 2},
					[]byte{0xff, 0xff, 0x00, 0x00, 0xcc, 0xcc},
				),
			},
			aez.CCDImage{PackData: aez.CCDImagePackData{NBC: 2}, BadColumns: []uint16{0xffff, 0x0000}},
			aez.SID(0),
			aez.CCD2,
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := make(chan common.DataRecord)
			target := make(chan common.DataRecord)
			go DecodeAEZ(target, source)
			source <- tt.arg
			close(source)
			got := <-target
			if (got.Error != nil) != tt.wantErr {
				t.Errorf("DataRecord.Error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want) {
				t.Errorf("DataRecord.Data = %v, want %v", got.Data, tt.want)
			}
			if got.SID != tt.wantSID {
				t.Errorf("DataRecord.SID = %v, want %v", got.SID, tt.wantSID)
			}
			if got.RID != tt.wantRID {
				t.Errorf("DataRecord.RID = %v, want %v", got.RID, tt.wantRID)
			}
			if len(got.Buffer) != tt.wantBufLen {
				t.Errorf("DataRecord.Buffer = %v, want length %v", got.Buffer, tt.wantBufLen)
			}
		})
	}
}

func Test_addData(t *testing.T) {
	type args struct {
		sourcePacket  common.DataRecord
		reader        *bytes.Reader
		exportable    common.Exportable
		exportableErr error
	}
	tests := []struct {
		name          string
		args          args
		wantData      common.Exportable
		wantBufferLen int
		wantErr       bool
	}{
		{
			"Adds the info to the DataRecord",
			args{
				common.DataRecord{},
				bytes.NewReader([]byte{0xff}),
				aez.STAT{},
				nil,
			},
			aez.STAT{},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addData(
				tt.args.sourcePacket,
				tt.args.reader,
				tt.args.exportable,
				tt.args.exportableErr,
			)
			if !reflect.DeepEqual(got.Data, tt.wantData) {
				t.Errorf("addData().Data = %v, want %v", got.Data, tt.wantData)
			}
			if len(got.Buffer) != tt.wantBufferLen {
				t.Errorf("addData().Buffer = %v, want length %v", got.Buffer, tt.wantBufferLen)
			}
			if (got.Error != nil) != tt.wantErr {
				t.Errorf("addData().Error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}

		})
	}
}
