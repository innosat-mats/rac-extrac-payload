package timeseries

import (
	"log"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

type CSVFactory func(pkg *common.DataRecord, stream OutStream) (CSVWriter, error)

// CSVCollection hold all active csv writers
type CSVCollection struct {
	streams map[OutStream]CSVWriter
	factory CSVFactory
}

func NewCollection(factory CSVFactory) CSVCollection {
	return CSVCollection{factory: factory, streams: make(map[OutStream]CSVWriter)}
}

// Write adds a csv row into the relevant out stream
func (collection *CSVCollection) Write(pkg *common.DataRecord) error {
	var writer CSVWriter
	var ok bool
	var err error
	stream := OutStreamFromDataRecord(pkg)
	if stream == Unknown {
		log.Printf("Unknown timeseries stream RID %v, SID %v", pkg.RID, pkg.SID)
		return nil
	}

	writer, ok = collection.streams[stream]
	if !ok {
		writer, err = collection.factory(pkg, stream)
		if err != nil {
			return err
		}
		if writer != nil {
			err := writer.SetSpecifications((*pkg).CSVSpecifications())
			if err != nil {
				return err
			}
			err = writer.SetHeaderRow((*pkg).CSVHeaders())
			if err != nil {
				return err
			}

		}
		collection.streams[stream] = writer
	}
	return writer.WriteData(pkg.CSVRow())
}

// CloseAll closes all open streams
func (collection *CSVCollection) CloseAll() {
	for stream := range collection.streams {
		oldWriter, ok := collection.streams[stream]
		if ok {
			oldWriter.Close()
			collection.streams[stream] = nil
		}
	}
}
