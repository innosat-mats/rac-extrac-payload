package exports

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/awstools"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/timeseries"
)

const awsS3Region = "eu-north-1"

func csvAWSWriterFactoryCreator(
	upload awstools.AWSUploadFunc,
	uploader *s3manager.Uploader,
	project string,
) timeseries.CSVFactory {
	return func(pkg *common.DataRecord, stream timeseries.OutStream) (timeseries.CSVWriter, error) {
		key := fmt.Sprintf("%v/%v.csv", project, stream.String())
		return timeseries.NewCSV(awstools.NewTimeseries(upload, uploader, key), key), nil
	}
}

// AWSS3CallbackFactory generates callbacks for writing to S3 instead of disk
func AWSS3CallbackFactory(
	upload awstools.AWSUploadFunc,
	project string,
	awsDescriptionPath string,
	writeImages bool,
	writeTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown) {

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(awsS3Region)}))
	uploader := s3manager.NewUploader(sess)
	timeseriesCollection := timeseries.NewCollection(
		csvAWSWriterFactoryCreator(upload, uploader, project),
	)

	if awsDescriptionPath != "" {
		awsDescription, err := os.Open(awsDescriptionPath)
		if err != nil {
			log.Fatalf("Could not find %v: %v", awsDescriptionPath, err)
		}
		key := fmt.Sprintf("ABOUT%v", filepath.Ext(awsDescriptionPath))
		if project != "" {
			key = fmt.Sprintf("%v/%v", project, key)
		}
		upload(uploader, key, awsDescription)
	}

	callback := func(pkg common.DataRecord) {
		if pkg.Error != nil {
			log.Println(pkg.Error)
			return
		}
		recoverWrite := func(imageFileName string) {
			if r := recover(); r != nil {
				log.Printf(
					"Could not process image %s (%v)",
					imageFileName, r,
				)
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
					imgFileName := ccdImage.FullImageName(project)
					defer recoverWrite(imgFileName)
					img := ccdImage.Image(pkg.Buffer)
					pngBuffer := bytes.NewBuffer([]byte{})
					err := png.Encode(pngBuffer, img)
					if err != nil {
						log.Panicf("failed encoding %s: %s", imgFileName, err)
					}
					upload(uploader, imgFileName, pngBuffer)
					jsonFileName := GetJSONFilename(imgFileName)
					jsonBuffer := bytes.NewBuffer([]byte{})
					WriteJSON(jsonBuffer, &pkg, jsonFileName)
					upload(uploader, jsonFileName, jsonBuffer)
				}()
			}
		}

		if writeTimeseries && pkg.Data != nil {
			// Write to the dedicated target stream
			err := timeseriesCollection.Write(&pkg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	teardown := func() {
		wg.Wait()
		timeseriesCollection.CloseAll()
	}

	return callback, teardown
}
