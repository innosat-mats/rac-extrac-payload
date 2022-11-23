package exports

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/timeseries"
)

func parquetName(dir string, pkg *common.DataRecord) string {
	name := timeseries.ParquetName(pkg)
	return filepath.Join(dir, name)
}

func parquetFileWriterFactoryCreator(
	dir string,
) timeseries.ParquetFactory {
	return func(pkg *common.DataRecord) (timeseries.ParquetWriter, error) {
		outPath := parquetName(dir, pkg)

		err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("could not create output prefix '%v'", outPath)
		}
		return timeseries.NewParquet(outPath, pkg), nil
	}
}

// DiskCallbackFactory returns a callback for disk writes
func ParquetCallbackFactory(
	output string,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown) {
	var err error
	timeseriesCollection := timeseries.NewParquetCollection(parquetFileWriterFactoryCreator(output))
	errorStats := common.NewErrorStats()

	// Create Directory and File
	err = os.MkdirAll(output, os.ModePerm)
	if err != nil {
		log.Printf("Could not create output directory '%v'", output)
	}

	callback := func(pkg common.DataRecord) {
		errorStats.Register(pkg.Error)
		if pkg.Error != nil {
			pkg.Error = fmt.Errorf(
				"%s %s",
				pkg.Error,
				common.MakePackageInfo(&pkg),
			)
			log.Println(pkg.Error)
		}

		if pkg.Data != nil {
			// Write to the dedicated target stream
			err = timeseriesCollection.Write(&pkg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	teardown := func() {
		timeseriesCollection.CloseAll()
		wg.Wait()
		log.Println(errorStats.Summarize())
	}

	return callback, teardown
}
