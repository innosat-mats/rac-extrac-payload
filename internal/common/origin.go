package common

import "time"

// OriginDescription describes the origin of ramses packages
type OriginDescription struct {
	Name string    // Name of the batch or file
	Date time.Time // Runtime of the batch
}
