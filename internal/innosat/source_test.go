package innosat

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVersion(t *testing.T) {
	sph := SourcePacketHeader{
		0b1110000000000000,
		0,
		0}

	ans := sph.PacketID.Version()
	if ans != 0b111 {
		t.Errorf("SourcePacketHeader.PacketID.Version() = %03b; want 111", ans)
	}
}

func TestType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0001000000000000,
		0,
		0}

	ans := sph.PacketID.Type()
	if ans != TC {
		t.Errorf("SourcePacketHeader.PacketID.Type() = %b; want 1", ans)
	}
}
func TestHeaderType(t *testing.T) {
	sph := SourcePacketHeader{
		0b0000100000000000,
		0,
		0}

	ans := sph.PacketID.HeaderType()
	if ans != 1 {
		t.Errorf("SourcePacketHeader.PacketID.HeaderType() = %b; want 1", ans)
	}
}

func TestGroupingFlags(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b1100000000000000,
		0}

	ans := sph.PacketSequenceControl.GroupingFlags()
	if ans != 0b11 {
		t.Errorf("SourcePacketHeader.PacketSequenceControl.GroupingFlags() = %02b; want 11", ans)
	}
}

func TestSequenceCount(t *testing.T) {
	sph := SourcePacketHeader{
		0,
		0b0011111111111111,
		0}

	ans := sph.PacketSequenceControl.SequenceCount()
	if ans != 0b11111111111111 {
		t.Errorf("SourcePacketHeader.PacketSequenceControl.SequenceCount() = %02b; want 11111111111111", ans)
	}
}

func TestSourcePacketHeader_APID(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   SourcePacketAPIDType
	}{
		{"Returns IdleAPID when apid is 0", fields{}, TimeAPID},
		{"Returns MainAPID when apid is 100", fields{PacketID: 2148}, MainAPID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			if got := sph.PacketID.APID(); got != tt.want {
				t.Errorf("SourcePacketHeader.PackedID.APID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePacketHeaderType_String(t *testing.T) {
	tests := []struct {
		name       string
		headerType SourcePacketHeaderType
		want       string
	}{
		{"TM is TM", TM, "TM"},
		{"TC is TC", TC, "TC"},
		{"5 is Unknown", SourcePacketHeaderType(5), "Unknown HeaderType 5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.headerType.String(); got != tt.want {
				t.Errorf("SourcePacketHeaderType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePackageContinuationFlagType_String(t *testing.T) {
	tests := []struct {
		name             string
		continuationFlag SourcePackageContinuationFlagType
		want             string
	}{
		{"SPCont", SPCont, "Continuation"},
		{"SPStart", SPStart, "Start"},
		{"SPCont", SPStop, "Stop"},
		{"SPCont", SPStandalone, "Standalone"},
		{"10 is Unkown", SourcePackageContinuationFlagType(10), "Unknown ContinuationFlag 10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.continuationFlag.String(); got != tt.want {
				t.Errorf("SourcePackageContinuationFlagType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePacketHeader_CSVSpecifications(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Creates spec", fields{}, []string{"INNOSAT", Specification}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			if got := sph.CSVSpecifications(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePacketHeader.CSVSpecifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePacketHeader_CSVHeaders(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Generates headers", fields{}, []string{"SPSequenceCount"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			if got := sph.CSVHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePacketHeader.CSVHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePacketHeader_CSVRow(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"Generates data",
			fields{PacketSequenceControl: PacketSequenceControl(0xc003)},
			[]string{"3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			if got := sph.CSVRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePacketHeader.CSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcePacketHeader_MarshalJSON(t *testing.T) {
	type fields struct {
		PacketID              packetID
		PacketSequenceControl PacketSequenceControl
		PacketLength          uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"Marshals into expected json",
			fields{PacketSequenceControl: 0xc008},
			[]byte(fmt.Sprintf("{\"specification\":\"%s\",\"spSequenceCount\":8}", Specification)),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sph := &SourcePacketHeader{
				PacketID:              tt.fields.PacketID,
				PacketSequenceControl: tt.fields.PacketSequenceControl,
				PacketLength:          tt.fields.PacketLength,
			}
			got, err := sph.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePacketHeader.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePacketHeader.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
