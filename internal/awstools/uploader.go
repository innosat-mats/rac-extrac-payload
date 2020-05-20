package awstools

import (
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const awsBucket = "mats-l0-artifacts"

// AWSUpload uploads file content to target bucket
func AWSUpload(uploader *s3manager.Uploader, key string, body io.Reader) {
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(strings.ReplaceAll(key, "\\", "/")),
		Body:   body,
	})
	if err != nil {
		log.Printf("Failed to upload file %v, %v", key, err)
	}

}

// AWSUploadFunc is the signature of an AWS upload function
type AWSUploadFunc func(uploader *s3manager.Uploader, key string, body io.Reader)
