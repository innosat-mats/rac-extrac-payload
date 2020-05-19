package awstools

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func TestNewTimeseries(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("localhost")}))
	upload := s3manager.NewUploader(sess)
	var idxUp = 0
	uploads := make(map[string]int)

	var uploader = func(uploader *s3manager.Uploader, key string, bodyBuffer io.Reader) {
		buf, _ := ioutil.ReadAll(bodyBuffer)
		uploads[key] = len(buf)
		idxUp++
	}

	ts := NewTimeseries(uploader, upload, "myfile")
	// Write several times possible
	for i := 0; i < 3; i++ {
		n, err := ts.Write([]byte("test"))
		if n != 4 || err != nil {
			t.Errorf("Timeseries.Write() = %v %v, want 4 <nil>", n, err)
			return
		}
		if idxUp != 0 {
			t.Errorf("Timesereis.Write() initiated an unexpected upload %v", uploads)
			return
		}
	}

	// Uploads file upon closing
	ts.Close()
	if idxUp != 1 {
		t.Errorf(
			"Timeseries.Close() didn't upload one file (%v files sent) with content %v",
			idxUp,
			uploads,
		)
	}
	l, ok := uploads["myfile"]
	if !ok || l != 4*3 {
		t.Errorf(
			"Timeseries.Close() didn't upload 'myfile' with 12 bytes (%v, %v)",
			l,
			ok,
		)
	}

	// Can't write after closing
	n, err := ts.Write([]byte("hello"))
	if n != 0 || err == nil {
		t.Errorf(
			"Timeseries.Write() should have written nothing and err after closing (%v, %v)",
			n,
			err,
		)
	}
}
