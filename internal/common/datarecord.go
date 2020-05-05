package common

import (
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// DataRecord holds the full decode from one or many Ramses packages
type DataRecord struct {
	Origin       OriginDescription          // Describes the origin of the data like filename or data batch name
	RamsesHeader ramses.Ramses              // Ramses header information
	RamsesSecure ramses.Secure              // Ramses secure header information
	SourceHeader innosat.SourcePacketHeader // Source header from the innosat platform
	TMHeader     innosat.TMDataFieldHeader  // Data header information
	SID          aez.SID                    // SID of the Data if any
	Data         Exportable                 // The data payload itself, HK report, jpeg image etc.
	Error        error                      // First propagated error from the decoding process
	Buffer       []byte                     // Currently unprocessed data (payload)
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
	headers = append(headers, record.SourceHeader.CSVHeaders()...)
	headers = append(headers, record.TMHeader.CSVHeaders()...)
	headers = append(headers, "SID")
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
	row = append(row, record.SourceHeader.CSVRow()...)
	row = append(row, record.TMHeader.CSVRow()...)
	row = append(row, record.SID.String())
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
