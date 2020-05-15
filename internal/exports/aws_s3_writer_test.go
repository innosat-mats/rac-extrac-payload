package exports

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func TestAWSS3CallbackFactory(t *testing.T) {
	type args struct {
		project             string
		descriptionFileName string
		descriptionFileBody []byte
		writeImages         bool
		writeTimeseries     bool
	}
	type upload struct {
		key     string
		bodyLen int
	}
	tests := []struct {
		name    string
		args    args
		record  common.DataRecord
		uploads []upload
	}{
		{
			"Uploads description file",
			args{
				project:             "myproj",
				descriptionFileName: "myfile.txt",
				descriptionFileBody: []byte("Hello"),
			},
			common.DataRecord{},
			[]upload{{key: filepath.Join("myproj", "ABOUT.txt"), bodyLen: 5}},
		},
		{
			"Uploads description file without project",
			args{
				project:             "",
				descriptionFileName: "myfile.md",
				descriptionFileBody: []byte("Hello"),
			},
			common.DataRecord{},
			[]upload{{key: "ABOUT.md", bodyLen: 5}},
		},
		{
			"Doesn't upload description file if name empty",
			args{
				project:             "myproj",
				descriptionFileName: "",
			},
			common.DataRecord{},
			[]upload{},
		},
		{
			"Uploads image",
			args{
				project:     "myproj",
				writeImages: true,
			},
			common.DataRecord{
				Origin: common.OriginDescription{Name: "MyRac.rac"},
				Data: aez.CCDImage{
					PackData: aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
				},
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			},
			[]upload{
				{
					key:     filepath.Join("myproj", "MyRac_5000000000.png"),
					bodyLen: 76, // 8 + header
				},
				{
					key:     filepath.Join("myproj", "MyRac_5000000000.json"),
					bodyLen: 853, // length of the json
				},
			},
		},
		{
			"Doesn't uploads image when told not to",
			args{
				project:     "myproj",
				writeImages: false,
			},
			common.DataRecord{
				Origin: common.OriginDescription{Name: "MyRac.rac"},
				Data: aez.CCDImage{
					PackData: aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
				},
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			},
			[]upload{},
		},
		{
			"Doesn't uploads errors",
			args{
				project:     "myproj",
				writeImages: true,
			},
			common.DataRecord{
				Origin: common.OriginDescription{Name: "MyRac.rac"},
				Data: aez.CCDImage{
					PackData: aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
				},
				Error:  errors.New("here be dragons"),
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			},
			[]upload{},
		},
		{
			"Uploads everything",
			args{
				project:             "myproj",
				descriptionFileName: "info.json",
				descriptionFileBody: []byte("[42,42]"),
				writeImages:         true,
			},
			common.DataRecord{
				Origin: common.OriginDescription{Name: "MyRac.rac"},
				Data: aez.CCDImage{
					PackData: aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
				},
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			},
			[]upload{
				{
					key:     filepath.Join("myproj", "ABOUT.json"),
					bodyLen: 7,
				},
				{
					key:     filepath.Join("myproj", "MyRac_5000000000.png"),
					bodyLen: 76, // 8 + header
				},
				{
					key:     filepath.Join("myproj", "MyRac_5000000000.json"),
					bodyLen: 853, // length of the json
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			var idxUp = 0

			var uploader = func(uploader *s3manager.Uploader, key string, bodyBuffer io.Reader) {
				buf, _ := ioutil.ReadAll(bodyBuffer)
				if idxUp >= len(tt.uploads) {
					t.Errorf(
						"Got unexpected upload #%v, key '%v', body %v",
						idxUp,
						key,
						buf,
					)
				} else {
					upload := tt.uploads[idxUp]
					if key != upload.key {
						t.Errorf("Upload %v: key = %v, want %v ", idxUp, key, upload.key)
					}
					if upload.bodyLen != len(buf) {
						t.Errorf("Upload %v: len(buf) = %v, want %v ", idxUp, len(buf), upload.bodyLen)
					}
				}

				idxUp++
			}

			awsDescriptionPath := ""
			if tt.args.descriptionFileName != "" {
				dir, err := ioutil.TempDir("", "innosat-mats")
				if err != nil {
					t.Errorf(
						"AWSS3CallbackFactory() could not setup test directory '%v': %v",
						dir,
						err,
					)
				}
				defer os.RemoveAll(dir)
				awsDescriptionPath = filepath.Join(dir, tt.args.descriptionFileName)
				f, err := os.Create(awsDescriptionPath)
				if err != nil {
					t.Errorf(
						"AWSS3CallbackFactory() could not setup test description file '%v': %v",
						awsDescriptionPath,
						err)
				}
				defer f.Close()
				n, err := f.Write(tt.args.descriptionFileBody)
				if n != len(tt.args.descriptionFileBody) || err != nil {
					t.Errorf(
						"AWSS3CallbackFactory() could not write test description file: %v, %v",
						n,
						err,
					)
				}
				f.Close()
			}

			callback, teardown := AWSS3CallbackFactory(
				uploader,
				tt.args.project,
				awsDescriptionPath,
				tt.args.writeImages,
				tt.args.writeTimeseries,
				&wg,
			)

			callback(tt.record)
			teardown()

			if idxUp < len(tt.uploads) {
				t.Errorf("Recorded %v uploads, want %v", idxUp, len(tt.uploads))
			}
		})
	}
}
