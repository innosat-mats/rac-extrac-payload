package awstools

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const awsBucket = "mats-l0-artifacts"

// AWSUpload uploads file content to target bucket
func AWSUpload(uploader *manager.Uploader, key string, body io.Reader) {
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(strings.ReplaceAll(key, "\\", "/")),
		Body:   body,
	})
	if err != nil {
		log.Printf("Failed to upload file %v, %v", key, err)
	}
}

// AWSUploadFunc is the signature of an AWS upload function
type AWSUploadFunc func(uploader *manager.Uploader, key string, body io.Reader)
