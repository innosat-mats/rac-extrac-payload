package exports

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func csvName(dir string, originName string, packetType string) string {
	name := fmt.Sprintf(
		"%v_%v.csv",
		strings.TrimSuffix(filepath.Base(originName), filepath.Ext(originName)),
		packetType,
	)
	return filepath.Join(dir, name)
}

func csvFileWriterFactory(
	dir string,
	originName string,
	packetType string,
	pkg *common.DataRecord,
) (CsvFileWriter, error) {
	outPath := csvName(dir, originName, packetType)

	out, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create output file '%v'", outPath)
	}

	// Make a csvFile and produce specs and header row
	csvFile := NewCSVFile(out, outPath)
	err = csvFile.SetSpecifications((*pkg).CSVSpecifications())
	if err != nil {
		return nil, err
	}
	err = csvFile.SetHeaderRow((*pkg).CSVHeaders())
	if err != nil {
		return nil, err
	}
	return &csvFile, nil
}

// DiskCallbackFactory returns a callback for disk writes
func DiskCallbackFactory(
	output string,
	writeImages bool,
	writeTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown) {
	var currentOrigin string = ""
	var err error
	var htrOut CsvFileWriter = nil
	var pwrOut CsvFileWriter = nil
	var cpruOut CsvFileWriter = nil
	var statOut CsvFileWriter = nil
	var pmOut CsvFileWriter = nil
	var ccdOut CsvFileWriter = nil
	var tcvOut CsvFileWriter = nil

	if writeImages || writeTimeseries {
		// Create Directory and File
		err := os.MkdirAll(output, os.ModePerm)
		if err != nil {
			log.Printf("Could not create output directory '%v'", output)
		}
	}

	callback := func(pkg common.DataRecord) {
		if pkg.Error != nil {
			log.Println(pkg.Error)
		}
		if writeImages {
			switch pkg.Data.(type) {
			case aez.CCDImage:
				ccdImage, ok := pkg.Data.(aez.CCDImage)
				if !ok {
					log.Print("Could not understand packet as CCDImage, this should be impossible.")
					break
				}

				wg.Add(1)
				go func() {
					defer wg.Done()

					img, imgFileName := ccdImage.Image(pkg.Buffer, output, pkg.Origin.Name)
					imgFile, err := os.Create(imgFileName)
					if err != nil {
						log.Printf("failed creating %s: %s", imgFileName, err)
						panic(err.Error())
					}
					defer imgFile.Close()
					png.Encode(imgFile, img)

					jsonFileName := GetJSONFilename(imgFileName)
					jsonFile, err := os.Create(jsonFileName)
					defer jsonFile.Close()
					if err != nil {
						log.Printf("failed creating %s: %s", jsonFileName, err)
						panic(err.Error())
					}
					WriteJSON(jsonFile, &pkg, jsonFileName)
				}()

			}
		}

		if !writeTimeseries {
			return
		}
		// Close streams from previous file
		if pkg.Origin.Name != currentOrigin {
			if pwrOut != nil {
				pwrOut.Close()
				pwrOut = nil
			}
			if htrOut != nil {
				htrOut.Close()
				htrOut = nil
			}
			if statOut != nil {
				statOut.Close()
				statOut = nil
			}
			if cpruOut != nil {
				cpruOut.Close()
				cpruOut = nil
			}
			if pmOut != nil {
				pmOut.Close()
				pmOut = nil
			}
			if ccdOut != nil {
				ccdOut.Close()
				ccdOut = nil
			}
			if tcvOut != nil {
				tcvOut.Close()
				tcvOut = nil
			}
			currentOrigin = pkg.Origin.Name
		}

		// We have nowhere to write partial extraction of record so we discard
		if pkg.Data == nil {
			return
		}

		// Write to the dedicated target stream
		switch pkg.Data.(type) {
		case aez.STAT:
			if statOut == nil {
				statOut, err = csvFileWriterFactory(output, currentOrigin, "STAT", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = statOut.WriteData(pkg.CSVRow())
		case aez.HTR:
			if htrOut == nil {
				htrOut, err = csvFileWriterFactory(output, currentOrigin, "HTR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = htrOut.WriteData(pkg.CSVRow())
		case aez.PWR:
			if pwrOut == nil {
				pwrOut, err = csvFileWriterFactory(output, currentOrigin, "PWR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = pwrOut.WriteData(pkg.CSVRow())
		case aez.CPRU:
			if cpruOut == nil {
				cpruOut, err = csvFileWriterFactory(output, currentOrigin, "CPRU", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = cpruOut.WriteData(pkg.CSVRow())
		case aez.PMData:
			if pmOut == nil {
				pmOut, err = csvFileWriterFactory(output, currentOrigin, "PM", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = pmOut.WriteData(pkg.CSVRow())
		case aez.CCDImage:
			if ccdOut == nil {
				ccdOut, err = csvFileWriterFactory(output, currentOrigin, "CCD", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = ccdOut.WriteData(pkg.CSVRow())
		case aez.TCAcceptSuccessData, aez.TCAcceptFailureData, aez.TCExecSuccessData, aez.TCExecFailureData:
			if tcvOut == nil {
				tcvOut, err = csvFileWriterFactory(output, currentOrigin, "TCV", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = tcvOut.WriteData(pkg.CSVRow())
		}
		// This error comes from writing a line and most probably would be a column missmatch
		// that means we should be able to continue and just report the error
		if err != nil {
			log.Println(err)
		}
	}

	teardown := func() {
		if statOut != nil {
			statOut.Close()
		}
		if htrOut != nil {
			htrOut.Close()
		}
		if pwrOut != nil {
			pwrOut.Close()
		}
		if cpruOut != nil {
			cpruOut.Close()
		}
		if pmOut != nil {
			pmOut.Close()
		}
		if ccdOut != nil {
			ccdOut.Close()
		}
		if tcvOut != nil {
			tcvOut.Close()
		}
		wg.Wait()
	}

	return callback, teardown
}
