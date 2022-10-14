package awstools

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

// Timeseries holds csv-buffer and upload features
type Timeseries struct {
	reader   io.Reader
	writer   io.Writer
	upload   AWSUploadFunc
	uploader *manager.Uploader
	key      string
}

func (ts *Timeseries) Write(data []byte) (int, error) {
	if ts.writer == nil {
		return 0, fmt.Errorf("Timeseries %v already closed", ts.key)
	}
	return ts.writer.Write(data)
}

// Close uploads data to aws
func (ts *Timeseries) Close() {
	ts.writer = nil
	ts.upload(ts.uploader, ts.key, ts.reader)
}

// NewTimeseries returns a timeseries that will invoke upload upon close
func NewTimeseries(upload AWSUploadFunc, uploader *manager.Uploader, key string) *Timeseries {
	buf := bytes.NewBuffer([]byte{})
	return &Timeseries{buf, buf, upload, uploader, key}
}
