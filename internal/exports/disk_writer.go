package exports

import (
	"encoding/csv"
	"encoding/json"
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

func csvOutputFactory(dir string, originName string, packetType string, pkg *common.DataRecord) (csvOutput, error) {
	outPath := csvName(dir, originName, packetType)

	out, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create output file '%v'", outPath)
	}

	// Make a csvFile and produce specs and header row
	csvFile := csvFile{File: out, Writer: csv.NewWriter(out)}
	err = csvFile.setSpecifications((*pkg).CSVSpecifications())
	if err != nil {
		return nil, err
	}
	err = csvFile.setHeaderRow((*pkg).CSVHeaders())
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
	var htrOut csvOutput = nil
	var pwrOut csvOutput = nil
	var cpruOut csvOutput = nil
	var statOut csvOutput = nil
	var pmOut csvOutput = nil
	var ccdOut csvOutput = nil
	var tcvOut csvOutput = nil

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
				imgFileName := getGrayscaleImageName(output, pkg.Origin.Name, ccdImage.PackData)
				wg.Add(1)
				go func() {
					imgData := getImageData(
						pkg.Buffer,
						ccdImage.PackData,
						imgFileName,
					)
					_, shift, _ := ccdImage.PackData.WDW.InputDataWindow()
					img := getGrayscaleImage(
						imgData,
						int(ccdImage.PackData.NCOL+aez.NCOLStartOffset),
						int(ccdImage.PackData.NROW),
						shift,
						imgFileName,
					)

					defer wg.Done()
					imgFile, err := os.Create(imgFileName)
					if err != nil {
						log.Printf("failed creating %s: %s", imgFileName, err)
						panic(err.Error())
					}
					defer imgFile.Close()
					png.Encode(imgFile, img)
				}()

				ext := filepath.Ext(imgFileName)
				jsonFileName := fmt.Sprintf(
					"%v.json",
					imgFileName[0:len(imgFileName)-len(ext)],
				)
				jsonFile, err := os.Create(jsonFileName)
				defer jsonFile.Close()
				if err != nil {
					log.Printf("failed creating %s: %s", jsonFileName, err)
					panic(err.Error())
				}
				err = json.NewEncoder(jsonFile).Encode(&pkg)
				if err != nil {
					log.Printf("failed to encode json into %s", jsonFileName)
				}
			}
		}

		if !writeTimeseries {
			return
		}
		// Close streams from previous file
		if pkg.Origin.Name != currentOrigin {
			if pwrOut != nil {
				pwrOut.close()
				pwrOut = nil
			}
			if htrOut != nil {
				htrOut.close()
				htrOut = nil
			}
			if statOut != nil {
				statOut.close()
				statOut = nil
			}
			if cpruOut != nil {
				cpruOut.close()
				cpruOut = nil
			}
			if pmOut != nil {
				pmOut.close()
				pmOut = nil
			}
			if ccdOut != nil {
				ccdOut.close()
				ccdOut = nil
			}
			if tcvOut != nil {
				tcvOut.close()
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
				statOut, err = csvOutputFactory(output, currentOrigin, "STAT", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = statOut.writeData(pkg.CSVRow())
		case aez.HTR:
			if htrOut == nil {
				htrOut, err = csvOutputFactory(output, currentOrigin, "HTR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = htrOut.writeData(pkg.CSVRow())
		case aez.PWR:
			if pwrOut == nil {
				pwrOut, err = csvOutputFactory(output, currentOrigin, "PWR", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = pwrOut.writeData(pkg.CSVRow())
		case aez.CPRU:
			if cpruOut == nil {
				cpruOut, err = csvOutputFactory(output, currentOrigin, "CPRU", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = cpruOut.writeData(pkg.CSVRow())
		case aez.PMData:
			if pmOut == nil {
				pmOut, err = csvOutputFactory(output, currentOrigin, "PM", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = pmOut.writeData(pkg.CSVRow())
		case aez.CCDImage:
			if ccdOut == nil {
				ccdOut, err = csvOutputFactory(output, currentOrigin, "CCD", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = ccdOut.writeData(pkg.CSVRow())
		case aez.TCAcceptSuccessData, aez.TCAcceptFailureData, aez.TCExecSuccessData, aez.TCExecFailureData:
			if tcvOut == nil {
				tcvOut, err = csvOutputFactory(output, currentOrigin, "TCV", &pkg)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = tcvOut.writeData(pkg.CSVRow())
		}
		// This error comes from writing a line and most probably would be a column missmatch
		// that means we should be able to continue and just report the error
		if err != nil {
			log.Println(err)
		}
	}

	teardown := func() {
		if statOut != nil {
			statOut.close()
		}
		if htrOut != nil {
			htrOut.close()
		}
		if pwrOut != nil {
			pwrOut.close()
		}
		if cpruOut != nil {
			cpruOut.close()
		}
		if pmOut != nil {
			pmOut.close()
		}
		if ccdOut != nil {
			ccdOut.close()
		}
		if tcvOut != nil {
			tcvOut.close()
		}
		wg.Wait()
	}

	return callback, teardown
}
