package common

import (
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// DataRecord holds the full decode from one or many Ramses packages
type DataRecord struct {
	Origin                   OriginDescription          `json:"origin"`                   // Describes the origin of the data like filename or data batch name
	RamsesHeader             ramses.Ramses              `json:"ramsesHeader"`             // Ramses header information
	OhbseCcsdsTMPacketHeader ramses.OhbseCcsdsTMPacket  `json:"ohbseCcsdsTMPacketHeader"` // The CCSDS compliant OHBSE TM Packet header
	SourceHeader             innosat.SourcePacketHeader `json:"sourceHeader"`             // Source header from the innosat platform
	TMHeader                 innosat.TMDataFieldHeader  `json:"tmHeader"`                 // Data header information
	SID                      aez.SID                    // SID of the Data if any
	RID                      aez.RID                    // RID of Data if any
	Data                     Exportable                 `json:"data"`            // The data payload itself, HK report, jpeg image etc.
	Error                    error                      `json:"error,omitempty"` // First propagated error from the decoding process
	Buffer                   []byte                     `json:"-"`               // Currently unprocessed data (payload)
}

// CSVSpecifications returns specifications used to generate content in CSV compatible format
func (record DataRecord) CSVSpecifications() []string {
	var specifications []string
	specifications = append(
		specifications,
		record.RamsesHeader.CSVSpecifications()...,
	)
	specifications = append(
		specifications,
		record.SourceHeader.CSVSpecifications()...,
	)

	if record.Data != nil {
		specifications = append(
			specifications,
			record.Data.CSVSpecifications()...,
		)
	}
	return specifications
}

// CSVHeaders returns a header row for the data record
func (record DataRecord) CSVHeaders() []string {
	var headers []string
	headers = append(headers, record.Origin.CSVHeaders()...)
	headers = append(headers, record.RamsesHeader.CSVHeaders()...)
	headers = append(headers, record.OhbseCcsdsTMPacketHeader.CSVHeaders()...)
	headers = append(headers, record.SourceHeader.CSVHeaders()...)
	headers = append(headers, record.TMHeader.CSVHeaders()...)
	headers = append(headers, "SID", "RID")
	if record.Data != nil {
		headers = append(headers, record.Data.CSVHeaders()...)
	}
	headers = append(headers, "Error")
	return headers
}

// CSVRow returns a data row for the record
func (record DataRecord) CSVRow() []string {
	var row []string
	row = append(row, record.Origin.CSVRow()...)
	row = append(row, record.RamsesHeader.CSVRow()...)
	row = append(row, record.OhbseCcsdsTMPacketHeader.CSVRow()...)
	row = append(row, record.SourceHeader.CSVRow()...)
	row = append(row, record.TMHeader.CSVRow()...)
	row = append(row, record.SID.String(), record.RID.String())
	if record.Data != nil {
		row = append(row, record.Data.CSVRow()...)
	}
	if record.Error != nil {
		row = append(row, record.Error.Error())
	} else {
		row = append(row, "")
	}
	return row
}

// AEZData returns the exportable aez data
func (record DataRecord) AEZData() interface{} {
	return record.Data
}

// OriginName returns the name of file/key that is source of data record
func (record DataRecord) OriginName() string {
	return record.Origin.Name
}

// RemainingBuffer returns the buffer in the packet not yet parsed
func (record DataRecord) RemainingBuffer() []byte {
	return record.Buffer
}

// ParsingError returns the error if any
func (record DataRecord) ParsingError() error {
	return record.Error
}
