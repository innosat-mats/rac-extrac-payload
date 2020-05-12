package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

/*
	TCAcceptSuccess:
*/

// TCAcceptSuccessData Telecommand Acceptance Report - Success
type TCAcceptSuccessData struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy if the Sequence Control Header field of the TC
}

// Read TCAcceptSuccess
func (tcv *TCAcceptSuccessData) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCAcceptSuccessData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCAcceptSuccessData) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCAcceptSuccessData) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

/*
	TCAcceptFailure:
*/

// TCAcceptFailureData Telecommand Acceptance Report - Failure
type TCAcceptFailureData struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy if the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
}

// Read TCAcceptFailure
func (tcv *TCAcceptFailureData) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCAcceptFailureData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCAcceptFailureData) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCAcceptFailureData) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

/*
	TCExecSuccess:
*/

// TCExecSuccessData Telecommand Execution Report - Success
type TCExecSuccessData struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy if the Sequence Control Header field of the TC
}

// Read TCExecSuccess
func (tcv *TCExecSuccessData) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCExecSuccessData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCExecSuccessData) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCExecSuccessData) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

/*
	TCExecFailure:
*/

// TCExecFailureData Telecommand Execution Report - Failure
type TCExecFailureData struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy if the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
	// TODO? There is also an optional variable length data field
}

// Read TCExecFailure
func (tcv *TCExecFailureData) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCExecFailureData) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCExecFailureData) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCExecFailureData) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}
