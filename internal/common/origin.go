package common

import (
	"path/filepath"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
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

// SetParquet sets the parquet representation of the OriginDescription
func (origin *OriginDescription) SetParquet(row *parquetrow.ParquetRow) {
	row.OriginFile = filepath.Base(origin.Name)
	row.ProcessingTime = origin.ProcessingDate
}
