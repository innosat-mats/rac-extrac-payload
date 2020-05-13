package common

//Exportable interface for data that can be written to target
type Exportable interface {
	CSVHeaders() []string
	CSVRow() []string
	CSVSpecifications() []string
}
