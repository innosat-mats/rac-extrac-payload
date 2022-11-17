package timeseries

import (
	"io"
	"log"
	"os"

	goparquet "github.com/fraugster/parquet-go"
	"github.com/fraugster/parquet-go/floor"
	"github.com/fraugster/parquet-go/parquetschema"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// Parquet gives easy access for parquet writing
type Parquet struct {
	writer        io.Writer
	parquetWriter *floor.Writer
	Name          string
	NHeaders      int
}

// NewParquet returns a Timeseries as parquet
func NewParquet(out io.Writer, name string, pkg *common.DataRecord) ParquetWriter {
	sd, err := parquetschema.ParseSchemaDefinition(RacSchema)
	if err != nil {
		log.Fatalf("could not parse parquet schema definition: %v", err)
	}
	metadata := pkg.ParquetSpecifications()
	writer, err := floor.NewFileWriter(
		name,
		goparquet.WithSchemaDefinition(sd),
		goparquet.WithMetaData(metadata),
	)
	if err != nil {
		log.Fatalf("could not create %v: %v", name, err)
	}
	return &Parquet{writer: out, parquetWriter: writer, Name: name}
}

// ParquetWriter implements ease of use writing functions
type ParquetWriter interface {
	Close()
	WriteData(data interface{}) error
}

// Close flushes and closes underlying file if any
func (parquet *Parquet) Close() {
	f, ok := parquet.writer.(*os.File)
	if ok {
		f.Close()
	}
	parquet.parquetWriter.Close()
}

// WriteData writes a data row
func (parquet *Parquet) WriteData(data interface{}) error {
	parquet.parquetWriter.Write(data)
	return nil
}
