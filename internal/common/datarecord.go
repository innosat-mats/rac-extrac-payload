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
	Origin         *OriginDescription          // Describes the origin of the data like filename or data batch name
	RamsesHeader   *ramses.Ramses              // Ramses header information
	RamsesTMHeader *ramses.TMHeader            // The CCSDS compliant OHBSE TM Packet header
	SourceHeader   *innosat.SourcePacketHeader // Source header from the innosat platform
	TMHeader       *innosat.TMHeader           // Data header information
	SID            aez.SID                     // SID of the Data if any
	RID            aez.RID                     // RID of Data if any
	Data           Exporter                    // The data payload itself, HK report, jpeg image etc.
	Error          error                       // First propagated error from the decoding process
	Buffer         []byte                      // Currently unprocessed data (payload)
}

// MarshalJSON makes a custom json of what is of interest in the struct
func (record *DataRecord) MarshalJSON() ([]byte, error) {
	var dataJSON []byte
	var dataJSONErr error
	switch record.Data.(type) {
	case *aez.CCDImage:
		ccd, ok := record.Data.(*aez.CCDImage)
		if ok {
			dataJSON, dataJSONErr = ccd.MarshalJSON()
		} else {
			dataJSON, dataJSONErr = json.Marshal("Could not marshal ccd data into json")
		}
	default:
		dataJSON, dataJSONErr = json.Marshal(record.Data)
	}
	buf, err := json.Marshal(&struct {
		Origin         *OriginDescription          `json:"origin"`
		RamsesHeader   *ramses.Ramses              `json:"ramsesHeader"`
		RamsesTMHeader *ramses.TMHeader            `json:"ramsesTMHeader"`
		SourceHeader   *innosat.SourcePacketHeader `json:"sourceHeader"`
		TMHeader       *innosat.TMHeader           `json:"tmHeader"`
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

	//Inject the specially marshalled Data in the right place
	jsonAsStr := string(buf)
	pattern := "\"data\":\"$1\""
	idx := strings.Index(jsonAsStr, pattern)
	return append(append(buf[0:idx+7], dataJSON...), buf[idx+len(pattern):]...), err
}

// CSVSpecifications returns specifications used to generate content in CSV compatible format
func (record *DataRecord) CSVSpecifications() []string {
	specifications := []string{"CODE", FullVersion()}
	if record.RamsesHeader != nil {
		specifications = append(
			specifications,
			record.RamsesHeader.CSVSpecifications()...,
		)
	} else {
		specifications = append(
			specifications,
			(&ramses.Ramses{}).CSVSpecifications()...,
		)
	}
	if record.SourceHeader != nil {
		specifications = append(
			specifications,
			record.SourceHeader.CSVSpecifications()...,
		)
	} else {
		specifications = append(
			specifications,
			(&innosat.SourcePacketHeader{}).CSVSpecifications()...,
		)
	}

	if record.Data != nil {
		specifications = append(
			specifications,
			record.Data.CSVSpecifications()...,
		)
	}
	return specifications
}

// CSVHeaders returns a header row for the data record
func (record *DataRecord) CSVHeaders() []string {
	var headers []string
	if record.Origin != nil {
		headers = append(headers, record.Origin.CSVHeaders()...)
	} else {
		headers = append(headers, (&OriginDescription{}).CSVHeaders()...)
	}
	if record.RamsesHeader != nil {
		headers = append(headers, record.RamsesHeader.CSVHeaders()...)
	} else {
		headers = append(headers, (&ramses.Ramses{}).CSVHeaders()...)
	}
	if record.RamsesTMHeader != nil {
		headers = append(headers, record.RamsesTMHeader.CSVHeaders()...)
	} else {
		headers = append(headers, (&ramses.TMHeader{}).CSVHeaders()...)
	}
	if record.SourceHeader != nil {
		headers = append(headers, record.SourceHeader.CSVHeaders()...)
	} else {
		headers = append(headers, (&innosat.SourcePacketHeader{}).CSVHeaders()...)
	}
	if record.TMHeader != nil {
		headers = append(headers, record.TMHeader.CSVHeaders()...)
	} else {
		headers = append(headers, (&innosat.TMHeader{}).CSVHeaders()...)
	}
	headers = append(headers, "SID", "RID")
	if record.Data != nil {
		headers = append(headers, record.Data.CSVHeaders()...)
	}
	headers = append(headers, "Error")
	return headers
}

// CSVRow returns a data row for the record
func (record *DataRecord) CSVRow() []string {
	var row []string
	if record.Origin != nil {
		row = append(row, record.Origin.CSVRow()...)
	} else {
		row = append(row, make([]string, len((&OriginDescription{}).CSVRow()))...)
	}
	if record.RamsesHeader != nil {
		row = append(row, record.RamsesHeader.CSVRow()...)
	} else {
		row = append(row, make([]string, len((&ramses.Ramses{}).CSVRow()))...)
	}
	if record.RamsesTMHeader != nil {
		row = append(row, record.RamsesTMHeader.CSVRow()...)
	} else {
		row = append(row, make([]string, len((&ramses.TMHeader{}).CSVRow()))...)
	}
	if record.SourceHeader != nil {
		row = append(row, record.SourceHeader.CSVRow()...)
	} else {
		row = append(row, make([]string, len((&innosat.SourcePacketHeader{}).CSVRow()))...)
	}
	if record.TMHeader != nil {
		row = append(row, record.TMHeader.CSVRow()...)
	} else {
		row = append(row, make([]string, len((&innosat.TMHeader{}).CSVRow()))...)
	}
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

// OriginName returns the origin name or empty string if unknown
func (record *DataRecord) OriginName() string {
	if record.Origin != nil {
		return record.Origin.Name
	}
	return ""
}
