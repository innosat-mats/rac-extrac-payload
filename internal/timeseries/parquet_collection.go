package timeseries

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// ParquetFactory is a function that creates ParquetWriters
type ParquetFactory func(pkg *common.DataRecord, stream OutStream) (ParquetWriter, error)

// ParquetCollection holds all active ParquetWriters
type ParquetCollection struct {
	streams map[string]ParquetWriter
	factory ParquetFactory
}

// NewParquetCollection returns a novel ready to use ParquetCollection
func NewParquetCollection(factory ParquetFactory) ParquetCollection {
	return ParquetCollection{factory: factory, streams: make(map[string]ParquetWriter)}
}

// ParquetName returns the whole name of the parquet, including partitioning prefix
func ParquetName(pkg *common.DataRecord, stream OutStream) string {
	tmTime := pkg.TMHeader.Time(time.Time{})
	prefix := filepath.Join(
		fmt.Sprintf("%v", tmTime.Year()),
		fmt.Sprintf("%v", int(tmTime.Month())),
		fmt.Sprintf("%v", tmTime.Day()),
	)
	baseName := filepath.Base(pkg.Origin.Name)
	ext := filepath.Ext(pkg.Origin.Name)
	name := fmt.Sprintf("%v_%v.parquet", stream.String(), strings.TrimSuffix(baseName, ext))
	return filepath.Join(prefix, name)
}

// Write adds a parquet row into the relevant out stream
func (collection ParquetCollection) Write(pkg *common.DataRecord) error {
	var writer ParquetWriter
	var ok bool
	var err error
	stream := OutStreamFromDataRecord(pkg)
	if stream == Unknown {
		log.Printf("Unknown timeseries stream RID %v, SID %v", pkg.RID, pkg.SID)
		return nil
	}
	streamName := ParquetName(pkg, stream)

	writer, ok = collection.streams[streamName]
	if !ok {
		writer, err = collection.factory(pkg, stream)
		if err != nil {
			return err
		}
		collection.streams[streamName] = writer
	}
	return writer.WriteData(GetParquetRow(pkg))
}

// CloseAll closes all open streams
func (collection *ParquetCollection) CloseAll() {
	for stream := range collection.streams {
		oldWriter, ok := collection.streams[stream]
		if ok {
			oldWriter.Close()
			collection.streams[stream] = nil
		}
	}
}
