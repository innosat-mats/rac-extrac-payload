package common

import (
	"path/filepath"
	"time"
)

// OriginDescription describes the origin of ramses packages
type OriginDescription struct {
	Name           string    `json:"name"`           // Name of the batch or file
	ProcessingDate time.Time `json:"processingTime"` // Runtime of the batch
}

// CSVHeaders returns the field names
func (origin *OriginDescription) CSVHeaders() []string {
	return []string{"OriginFile", "ProcessingDate"}
}

// CSVRow returns the field values
func (origin *OriginDescription) CSVRow() []string {
	return []string{
		filepath.Base(origin.Name),
		origin.ProcessingDate.Format(time.RFC3339),
	}
}

// OriginDescriptionParquet holds the parquet representation of the OriginDescription
type OriginDescriptionParquet struct {
	OriginFile     string    `parquet:"OriginFile"`
	ProcessingDate time.Time `parquet:"ProcessingTime"`
}

// GetParquet returns the parquet representation of the OriginDescription
func (origin *OriginDescription) GetParquet() OriginDescriptionParquet {
	return OriginDescriptionParquet{origin.Name, origin.ProcessingDate}
}
