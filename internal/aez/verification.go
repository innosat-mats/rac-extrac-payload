package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type TCVParquet struct {
	TCV       string `parquet:"TCV"`
	TCPID     uint16 `parquet:"TCPID"`
	PSC       uint16 `parquet:"PSC"`
	ErrorCode uint8  `parquet:"ErrorCode"`
}

/*
TCAcceptSuccess:
*/

// TCAcceptSuccessData Telecommand Acceptance Report - Success
type TCAcceptSuccessData struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy of the Sequence Control Header field of the TC
}

// NewTCAcceptSuccessData reads a TCAcceptSuccessData from buffer
func NewTCAcceptSuccessData(buf io.Reader) (*TCAcceptSuccessData, error) {
	tcv := TCAcceptSuccessData{}
	err := binary.Read(buf, binary.LittleEndian, &tcv)
	return &tcv, err
}

// CSVSpecifications returns the version of the spec used
func (tcv *TCAcceptSuccessData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv *TCAcceptSuccessData) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "TCV")
	headers = append(headers, csvHeader(tcv)...)
	return append(headers, "ErrorCode")
}

// CSVRow returns the data row
func (tcv *TCAcceptSuccessData) CSVRow() []string {
	var row []string
	row = append(row, "Accept")
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return append(row, "")
}

// GetParquet returns the parquet representation of the TCV
func (tcv *TCAcceptSuccessData) GetParquet() TCVParquet {
	return TCVParquet{
		TCV:   "Accept",
		PSC:   tcv.PSC,
		TCPID: tcv.TCPID,
	}
}

/*
	TCAcceptFailure:
*/

// TCAcceptFailureData Telecommand Acceptance Report - Failure
type TCAcceptFailureData struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy of the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
}

// NewTCAcceptFailureData reads a TCAcceptFailureData from buffer
func NewTCAcceptFailureData(buf io.Reader) (*TCAcceptFailureData, error) {
	tcv := TCAcceptFailureData{}
	err := binary.Read(buf, binary.LittleEndian, &tcv)
	return &tcv, err
}

// CSVSpecifications returns the version of the spec used
func (tcv *TCAcceptFailureData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv *TCAcceptFailureData) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "TCV")
	return append(headers, csvHeader(tcv)...)
}

// CSVRow returns the data row
func (tcv *TCAcceptFailureData) CSVRow() []string {
	var row []string
	row = append(row, "Accept")
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

// GetParquet returns the parquet representation of the TCV
func (tcv *TCAcceptFailureData) GetParquet() TCVParquet {
	return TCVParquet{
		TCV:       "Accept",
		PSC:       tcv.PSC,
		TCPID:     tcv.TCPID,
		ErrorCode: tcv.ErrorCode,
	}
}

/*
	TCExecSuccess:
*/

// TCExecSuccessData Telecommand Execution Report - Success
type TCExecSuccessData struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy of the Sequence Control Header field of the TC
}

// NewTCExecSuccessData reads a TCExecSuccessData from buffer
func NewTCExecSuccessData(buf io.Reader) (*TCExecSuccessData, error) {
	tcv := TCExecSuccessData{}
	err := binary.Read(buf, binary.LittleEndian, &tcv)
	return &tcv, err
}

// CSVSpecifications returns the version of the spec used
func (tcv *TCExecSuccessData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv *TCExecSuccessData) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "TCV")
	headers = append(headers, csvHeader(tcv)...)
	return append(headers, "ErrorCode")
}

// CSVRow returns the data row
func (tcv *TCExecSuccessData) CSVRow() []string {
	var row []string
	row = append(row, "Exec")
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return append(row, "")
}

// GetParquet returns the parquet representation of the TCV
func (tcv *TCExecSuccessData) GetParquet() TCVParquet {
	return TCVParquet{
		TCV:   "Accept",
		PSC:   tcv.PSC,
		TCPID: tcv.TCPID,
	}
}

/*
	TCExecFailure:
*/

// TCExecFailureData Telecommand Execution Report - Failure
type TCExecFailureData struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy of the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
	// TODO? There is also an optional variable length data field
}

// NewTCExecFailureData reads a TCExecFailureData from buffer
func NewTCExecFailureData(buf io.Reader) (*TCExecFailureData, error) {
	tcv := TCExecFailureData{}
	err := binary.Read(buf, binary.LittleEndian, &tcv)
	return &tcv, err
}

// CSVSpecifications returns the version of the spec used
func (tcv *TCExecFailureData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv *TCExecFailureData) CSVHeaders() []string {
	var headers []string
	headers = append(headers, "TCV")
	return append(headers, csvHeader(tcv)...)
}

// CSVRow returns the data row
func (tcv *TCExecFailureData) CSVRow() []string {
	var row []string
	row = append(row, "Exec")
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

// GetParquet returns the parquet representation of the TCV
func (tcv *TCExecFailureData) GetParquet() TCVParquet {
	return TCVParquet{
		TCV:       "Exec",
		PSC:       tcv.PSC,
		TCPID:     tcv.TCPID,
		ErrorCode: tcv.ErrorCode,
	}
}
