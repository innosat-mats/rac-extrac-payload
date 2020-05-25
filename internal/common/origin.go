package common

import "time"

// OriginDescription describes the origin of ramses packages
type OriginDescription struct {
	Name           string    `json:"name"`           // Name of the batch or file
	ProcessingDate time.Time `json:"processingTime"` // Runtime of the batch
}

//CSVHeaders returns the field names
func (origin *OriginDescription) CSVHeaders() []string {
	return []string{"File", "ProcessingDate"}
}

//CSVRow returns the field values
func (origin *OriginDescription) CSVRow() []string {
	return []string{
		origin.Name,
		origin.ProcessingDate.Format(time.RFC3339),
	}
}
