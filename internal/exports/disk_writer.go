package exports

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/timeseries"
)

func csvName(dir string, packetType string) string {
	name := fmt.Sprintf("%v.csv", packetType)
	return filepath.Join(dir, name)
}

func csvFileWriterFactoryCreator(
	dir string,
) timeseries.CSVFactory {
	return func(pkg *common.DataRecord, stream timeseries.OutStream) (timeseries.CSVWriter, error) {
		outPath := csvName(dir, stream.String())

		out, err := os.Create(outPath)
		if err != nil {
			return nil, fmt.Errorf("could not create output file '%v'", outPath)
		}
		return timeseries.NewCSV(out, outPath), nil
	}
}

// DiskCallbackFactory returns a callback for disk writes
func DiskCallbackFactory(
	output string,
	writeImages bool,
	writeTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown) {
	var err error
	timeseriesCollection := timeseries.NewCollection(csvFileWriterFactoryCreator(output))
	errorStats := common.NewErrorStats()

	if writeImages || writeTimeseries {
		// Create Directory and File
		err := os.MkdirAll(output, os.ModePerm)
		if err != nil {
			log.Printf("Could not create output directory '%v'", output)
		}
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
		recoverWrite := func(imageFileName string) {
			if r := recover(); r != nil {
				log.Printf(
					"Processing incomplete for image %s, skipping (%v)",
					imageFileName, r,
				)
				os.Remove(imageFileName)
				os.Remove(GetJSONFilename(imageFileName))
			}
		}
		if writeImages {
			switch pkg.Data.(type) {
			case *aez.CCDImage:
				ccdImage, ok := pkg.Data.(*aez.CCDImage)
				if !ok {
					log.Print("Could not understand packet as CCDImage, this should be impossible.")
					break
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					imgFileName := ccdImage.FullImageName(output)
					defer recoverWrite(imgFileName)
					img := ccdImage.Image(pkg.Buffer)
					imgFile, err := os.Create(imgFileName)
					if err != nil {
						log.Panicf("failed creating %s: %s", imgFileName, err)
					}
					defer imgFile.Close()
					err = png.Encode(imgFile, img)
					if err != nil {
						log.Panicf("failed encoding %s: %s", imgFileName, err)
					}
					jsonFileName := GetJSONFilename(imgFileName)
					jsonFile, err := os.Create(jsonFileName)
					if err != nil {
						log.Panicf("failed creating %s: %s", jsonFileName, err)
					}
					defer jsonFile.Close()
					WriteJSON(jsonFile, &pkg, jsonFileName)
				}()

			}
		}

		if writeTimeseries && pkg.Data != nil {
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
