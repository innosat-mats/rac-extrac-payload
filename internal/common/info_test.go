package common

import (
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func Test_MakePackageInfo(t *testing.T) {
	tests := []struct {
		name         string
		sourcePacket *DataRecord
		want         string
	}{
		{
			"Empty",
			&DataRecord{},
			"[]",
		},
		{
			"With name",
			&DataRecord{
				Origin: &OriginDescription{
					Name: "test.rac",
				},
			},
			"[test.rac]",
		},
		{
			"With name & SourceHeader",
			&DataRecord{
				Origin: &OriginDescription{
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
			&DataRecord{
				Origin: &OriginDescription{
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
			&DataRecord{
				Origin: &OriginDescription{
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
			if got := MakePackageInfo(tt.sourcePacket); got != tt.want {
				t.Errorf("MakePackageInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
