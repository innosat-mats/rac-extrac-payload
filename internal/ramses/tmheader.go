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

// TMHeader is the OHBSE CCSDS TM Packet Header in the specification
type TMHeader struct {
	_                [8]byte
	QualityIndicator QualityIndicator `json:"qualityIndicator"` // QualityIndicator indicates whether the transported data is complete or partial (0 = Complete, 1 = partial)
	LossFlag         LossFlag         `json:"lossFlag"`         // LossFlag is used to indicate that a sequence discontinuity has been detected
	VCFrameCounter   uint8            `json:"vcFrameCounter"`   // VCFrameCounter is a counter of the transfer fram the payload packet arrived in
	_                [5]byte
}

// NewTMHeader reads a TMHeader from buffer
func NewTMHeader(buf io.Reader) (*TMHeader, error) {
	header := TMHeader{}
	err := binary.Read(buf, binary.LittleEndian, &header)
	return &header, err
}

// CSVHeaders returns the field names
func (header *TMHeader) CSVHeaders() []string {
	return []string{
		"QualityIndicator",
		"LossFlag",
		"VCFrameCounter",
	}
}

// CSVRow returns the field values
func (header *TMHeader) CSVRow() []string {
	return []string{
		strconv.Itoa(int(header.QualityIndicator)),
		strconv.Itoa(int(header.LossFlag)),
		strconv.Itoa(int(header.VCFrameCounter)),
	}
}

// RamsesTMHeaderParquet holds the parquet representation of the TMHeader
type RamsesTMHeaderParquet struct {
	QualityIndicator uint8 `parquet:"QualityIndicator"`
	LossFlag         uint8 `parquet:"LossFlag"`
	VCFrameCounter   uint8 `parquet:"VCFrameCounter"`
}

// GetParquet returns the parquet representation of the TMHeader
func (header *TMHeader) GetParquet() RamsesTMHeaderParquet {
	return RamsesTMHeaderParquet{
		uint8(header.QualityIndicator),
		uint8(header.LossFlag),
		header.VCFrameCounter,
	}
}
