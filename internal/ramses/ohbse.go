package ramses

import (
	"encoding/binary"
	"io"
	"strconv"
)

// QualityIndicator indicates whether the transported data is complete or partial
type QualityIndicator uint8

const (
	// CompletePacket ...
	CompletePacket QualityIndicator = iota
	// IncompletePacket ...
	IncompletePacket
)

// LossFlag is used to indicate that a sqeunce discontinuity has been detected and that one or more dataunits may have been lost
type LossFlag uint8

const (
	// NoDiscontinuities detected
	NoDiscontinuities LossFlag = iota
	// Discontinuities detected
	Discontinuities
)

// OhbseCcsdsTMPacket is used to transport telemetry packets on CCSDS format
type OhbseCcsdsTMPacket struct {
	_                [8]byte
	QualityIndicator QualityIndicator `json:"qualityIndicator"` // QualityIndicator indicates whether the transported data is complete or partial (0 = Complete, 1 = partial)
	LossFlag         LossFlag         `json:"lossFlag"`         // LossFlag is used to indicate that a sequence discontinuity has been detected
	VCFrameCounter   uint8            `json:"vcFrameCounter"`   // VCFrameCounter is a counter of the transfer fram the payload packet arrived in
	_                [5]byte
}

// Read a OhbseCssdsTMPacket header
func (header *OhbseCcsdsTMPacket) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, header)
}

//CSVHeaders returns the field names
func (header OhbseCcsdsTMPacket) CSVHeaders() []string {
	return []string{
		"QualityIndicator",
		"LossFlag",
		"VCFrameCounter",
	}
}

//CSVRow returns the field values
func (header OhbseCcsdsTMPacket) CSVRow() []string {
	return []string{
		strconv.Itoa(int(header.QualityIndicator)),
		strconv.Itoa(int(header.LossFlag)),
		strconv.Itoa(int(header.VCFrameCounter)),
	}
}
