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

// TCAcceptSuccess Telecommand Acceptance Report - Success
type TCAcceptSuccess struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy if the Sequence Control Header field of the TC
}

// Read TCAcceptSuccess
func (tcv *TCAcceptSuccess) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCAcceptSuccess) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCAcceptSuccess) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCAcceptSuccess) CSVRow() []string {
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

// TCAcceptFailure Telecommand Acceptance Report - Failure
type TCAcceptFailure struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy if the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
}

// Read TCAcceptFailure
func (tcv *TCAcceptFailure) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCAcceptFailure) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCAcceptFailure) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCAcceptFailure) CSVRow() []string {
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

// TCExecSuccess Telecommand Execution Report - Success
type TCExecSuccess struct {
	TCPID uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC   uint16 // PSC is a copy if the Sequence Control Header field of the TC
}

// Read TCExecSuccess
func (tcv *TCExecSuccess) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCExecSuccess) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCExecSuccess) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCExecSuccess) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}

// TCExecFailure Telecommand Execution Report - Failure
type TCExecFailure struct {
	TCPID     uint16 // TCPID is a copy of the Packet ID header field of the TC
	PSC       uint16 // PSC is a copy if the Sequence Control Header field of the TC
	ErrorCode uint8  // Error code
	// TODO? There is also an optional variable length data field
}

/*
	TCExecFailure:
*/

// Read TCExecFailure
func (tcv *TCExecFailure) Read(buf io.Reader) error {
	return binary.Read(buf, binary.LittleEndian, tcv)
}

// CSVSpecifications returns the version of the spec used
func (tcv TCExecFailure) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the header row
func (tcv TCExecFailure) CSVHeaders() []string {
	return csvHeader(tcv)
}

// CSVRow returns the data row
func (tcv TCExecFailure) CSVRow() []string {
	var row []string
	val := reflect.Indirect(reflect.ValueOf(tcv))
	for i := 0; i < val.NumField(); i++ {
		row = append(row, fmt.Sprintf("%v", val.Field(i).Uint()))
	}
	return row
}
