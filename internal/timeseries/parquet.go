package timeseries

import (
	"log"

	goparquet "github.com/fraugster/parquet-go"
	"github.com/fraugster/parquet-go/floor"
	"github.com/fraugster/parquet-go/parquet"
	"github.com/fraugster/parquet-go/parquetschema"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
)

// Parquet gives easy access for parquet writing
type Parquet struct {
	parquetWriter *floor.Writer
	Name          string
	NHeaders      int
}

var streamToScheme = map[OutStream]string{
	Unknown: parquetrow.RacSchema,
	HTR:     parquetrow.RacHTRSchema,
	PWR:     parquetrow.RacPWRSchema,
	CPRU:    parquetrow.RacCPRUSchema,
	STAT:    parquetrow.RacSTATSchema,
	PM:      parquetrow.RacPMSchema,
	CCD:     parquetrow.RacCCDSchema,
	TCV:     parquetrow.RacTCVSchema,
}

// NewParquet returns a Timeseries as parquet
func NewParquet(name string, pkg *common.DataRecord) ParquetWriter {
	schema := streamToScheme[OutStreamFromDataRecord(pkg)]
	sd, err := parquetschema.ParseSchemaDefinition(schema)
	if err != nil {
		log.Fatalf("could not parse parquet schema definition: %v", err)
	}
	metadata := pkg.ParquetSpecifications()
	writer, err := floor.NewFileWriter(
		name,
		goparquet.WithSchemaDefinition(sd),
		goparquet.WithMetaData(metadata),
		goparquet.WithCompressionCodec(parquet.CompressionCodec_SNAPPY),
	)
	if err != nil {
		log.Fatalf("could not create %v: %v", name, err)
	}
	return &Parquet{parquetWriter: writer, Name: name}
}

// ParquetWriter implements ease of use writing functions
type ParquetWriter interface {
	Close()
	WriteData(data interface{}) error
}

// Close flushes and closes underlying file if any
func (parquet *Parquet) Close() {
	parquet.parquetWriter.Close()
}

// WriteData writes a data row
func (parquet *Parquet) WriteData(data interface{}) error {
	return parquet.parquetWriter.Write(data)
}

// GetParquetRow returns the exportable parquet representation of a record, including common attributes
func GetParquetRow(pkg *common.DataRecord) parquetrow.ParquetRow {
	row := parquetrow.ParquetRow{}
	pkg.SetParquet(&row)

	switch pkg.Data.(type) {
	case *aez.CCDImage:
		ccd, ok := pkg.Data.(*aez.CCDImage)
		if ok {
			ccd.SetParquet(&row, pkg.Buffer)
		}
	case *aez.PMData:
		pm, ok := pkg.Data.(*aez.PMData)
		if ok {
			pm.SetParquet(&row)
		}
	case *aez.HTR:
		htr, ok := pkg.Data.(*aez.HTR)
		if ok {
			htr.SetParquet(&row)
		}
	case *aez.PWR:
		pwr, ok := pkg.Data.(*aez.PWR)
		if ok {
			pwr.SetParquet(&row)
		}
	case *aez.CPRU:
		cpru, ok := pkg.Data.(*aez.CPRU)
		if ok {
			cpru.SetParquet(&row)
		}
	case *aez.STAT:
		stat, ok := pkg.Data.(*aez.STAT)
		if ok {
			stat.SetParquet(&row)
		}
	case *aez.TCAcceptSuccessData:
		tcv, ok := pkg.Data.(*aez.TCAcceptSuccessData)
		if ok {
			tcv.SetParquet(&row)
		}
	case *aez.TCAcceptFailureData:
		tcv, ok := pkg.Data.(*aez.TCAcceptFailureData)
		if ok {
			tcv.SetParquet(&row)
		}
	case *aez.TCExecSuccessData:
		tcv, ok := pkg.Data.(*aez.TCExecSuccessData)
		if ok {
			tcv.SetParquet(&row)
		}
	case *aez.TCExecFailureData:
		tcv, ok := pkg.Data.(*aez.TCExecFailureData)
		if ok {
			tcv.SetParquet(&row)
		}
	}
	return row
}
