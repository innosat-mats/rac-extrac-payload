package exports

//Exportable interface for data that can be written to target
type Exportable interface {
	CSVHeaders() []string
	CSVRow() []string
	CSVSpecifications() []string
}

//ExportablePackage interface for root packages that can be written to target
type ExportablePackage interface {
	CSVHeaders() []string
	CSVRow() []string
	CSVSpecifications() []string
	OriginName() string
	AEZData() interface{}
}
