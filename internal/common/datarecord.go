package common

import (
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// DataRecord holds the full decode from one or many Ramses packages
type DataRecord struct {
	Origin       OriginDescription          // Describes the origin of the data like filename or data batch name
	RamsesHeader ramses.Ramses              // Ramses header information
	RamsesSecure ramses.Secure              // Ramses secure header information
	SourceHeader innosat.SourcePacketHeader // Source header from the innosat platform
	TMHeader     innosat.TMDataFieldHeader  // Data header information
	Data         PackageType                // The data payload itself, HK report, jpeg image etc.
	Error        error                      // First propagated error from the decoding process
	Buffer       []byte                     // Currently unprocessed data (payload)
}
