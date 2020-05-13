package aez

import "fmt"

// CCDImage is a container for the invariant CCDImagePackData header and the variable BadColumns that follow
type CCDImage struct {
	PackData   CCDImagePackData
	BadColumns []uint16
}

// CSVSpecifications returns the specs used in creating the struct
func (ccd CCDImage) CSVSpecifications() []string {
	return []string{"Specification", Specification}
}

// CSVHeaders returns the exportable field names
func (ccd CCDImage) CSVHeaders() []string {
	return append(ccd.PackData.CSVHeaders(), "BC")
}

// CSVRow returns the exportable field values
func (ccd CCDImage) CSVRow() []string {
	row := ccd.PackData.CSVRow()
	return append(row, fmt.Sprintf("%v", ccd.BadColumns))
}
