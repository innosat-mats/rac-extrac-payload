package common

//Exporter interface for data that can be written to target
type Exporter interface {
	CSVHeaders() []string
	CSVRow() []string
	CSVSpecifications() []string
}
