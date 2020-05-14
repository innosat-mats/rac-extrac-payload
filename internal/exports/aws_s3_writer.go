package exports

import (
	"bytes"
	"image/png"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

const awsBucket = "mats-l0-artifacts"
const awsS3Region = "eu-north-1"

func upload(uploader *s3manager.Uploader, key string, pngBuffer io.Reader) {
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(strings.ReplaceAll(key, "\\", "/")),
		Body:   pngBuffer,
	})
	if err != nil {
		log.Printf("Failed to upload file %v, %v", key, err)
	}

}

// AWSS3CallbackFactory generates callbacks for writing to S3 instead of disk
func AWSS3CallbackFactory(
	project string,
	writeImages bool,
	writeTimeseries bool,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown) {

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(awsS3Region)}))

	uploader := s3manager.NewUploader(sess)

	callback := func(pkg common.DataRecord) {
		if pkg.Error != nil {
			log.Println(pkg.Error)
			return
		}
		if writeImages {
			switch pkg.Data.(type) {
			case aez.CCDImage:
				ccdImage, ok := pkg.Data.(aez.CCDImage)
				if !ok {
					log.Print("Could not understand packet as CCDImage, this should be impossible.")
					break
				}
				img, imgFileName := ccdImage.Image(pkg.Buffer, project, pkg.Origin.Name)
				wg.Add(1)
				go func() {
					pngBuffer := bytes.NewBuffer([]byte{})
					png.Encode(pngBuffer, img)
					upload(uploader, imgFileName, pngBuffer)

					jsonFileName := GetJSONFilename(imgFileName)
					jsonBuffer := bytes.NewBuffer([]byte{})
					WriteJSON(jsonBuffer, &pkg, jsonFileName)
					upload(uploader, jsonFileName, jsonBuffer)

					wg.Done()
				}()
			}
		}
	}

	teardown := func() {
		wg.Wait()
	}

	return callback, teardown
}
