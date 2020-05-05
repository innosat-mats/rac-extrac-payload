package aez

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
	return []string{}
}

// CSVRow returns the exportable field values
func (ccd CCDImage) CSVRow() []string {
	return []string{}
}
