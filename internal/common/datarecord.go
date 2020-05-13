package common

import (
	"encoding/json"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// DataRecord holds the full decode from one or many Ramses packages
type DataRecord struct {
	Origin         OriginDescription          // Describes the origin of the data like filename or data batch name
	RamsesHeader   ramses.Ramses              // Ramses header information
	RamsesTMHeader ramses.TMHeader            // The CCSDS compliant OHBSE TM Packet header
	SourceHeader   innosat.SourcePacketHeader // Source header from the innosat platform
	TMHeader       innosat.TMHeader           // Data header information
	SID            aez.SID                    // SID of the Data if any
	RID            aez.RID                    // RID of Data if any
	Data           Exportable                 // The data payload itself, HK report, jpeg image etc.
	Error          error                      // First propagated error from the decoding process
	Buffer         []byte                     // Currently unprocessed data (payload)
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (record *DataRecord) MarshalJSON() ([]byte, error) {
	var dataJSON []byte
	var dataJSONErr error
	switch record.Data.(type) {
	case aez.CCDImage:
		ccd, ok := record.Data.(aez.CCDImage)
		if ok {
			dataJSON, dataJSONErr = ccd.MarshalJSON()
		} else {
			dataJSON, dataJSONErr = json.Marshal("Could not marshal ccd data into json")
		}
	default:
		dataJSON, dataJSONErr = json.Marshal(record.Data)
	}
	buf, err := json.Marshal(&struct {
		Origin         OriginDescription          `json:"origin"`
		RamsesHeader   ramses.Ramses              `json:"ramsesHeader"`
		RamsesTMHeader ramses.TMHeader            `json:"ramsesTMHeader"`
		SourceHeader   innosat.SourcePacketHeader `json:"sourceHeader"`
		TMHeader       innosat.TMHeader           `json:"tmHeader"`
		SID            aez.SID
		RID            aez.RID
		Data           string `json:"data"`
		Error          error  `json:"error,omitempty"`
		JSONError      error  `json:"errorJSON,omitempty"`
	}{
		Origin:         record.Origin,
		RamsesHeader:   record.RamsesHeader,
		RamsesTMHeader: record.RamsesTMHeader,
		SourceHeader:   record.SourceHeader,
		TMHeader:       record.TMHeader,
		SID:            record.SID,
		RID:            record.RID,
		Data:           "$1",
		Error:          record.Error,
		JSONError:      dataJSONErr,
	})
	s := string(buf)
	pattern := "\"data\":\"$1\""
	idx := strings.Index(s, pattern)
	return append(append(buf[0:idx+7], dataJSON...), buf[idx+len(pattern):]...), err
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
	headers = append(headers, record.RamsesTMHeader.CSVHeaders()...)
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
	row = append(row, record.RamsesTMHeader.CSVRow()...)
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
