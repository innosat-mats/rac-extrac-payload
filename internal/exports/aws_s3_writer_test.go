package exports

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
	"github.com/innosat-mats/rac-extract-payload/internal/timeseries"
)

func TestAWSS3CallbackFactory(t *testing.T) {
	type args struct {
		project             string
		descriptionFileName string
		descriptionFileBody []byte
		writeImages         bool
		writeTimeseries     bool
	}
	tests := []struct {
		name    string
		args    args
		records []common.DataRecord
		uploads map[string]int
	}{
		{
			"Uploads description file",
			args{
				project:             "myproj",
				descriptionFileName: "myfile.txt",
				descriptionFileBody: []byte("Hello"),
			},
			[]common.DataRecord{{}},
			map[string]int{"myproj/ABOUT.txt": 5},
		},
		{
			"Uploads description file without project",
			args{
				project:             "",
				descriptionFileName: "myfile.md",
				descriptionFileBody: []byte("Hello"),
			},
			[]common.DataRecord{{}},
			map[string]int{"ABOUT.md": 5},
		},
		{
			"Doesn't upload description file if name empty",
			args{
				project:             "myproj",
				descriptionFileName: "",
			},
			[]common.DataRecord{{}},
			map[string]int{},
		},
		{
			"Doesn't upload timeseries for no data",
			args{
				project:             "myproj",
				descriptionFileName: "",
				writeTimeseries:     true,
			},
			[]common.DataRecord{{}},
			map[string]int{},
		},
		{
			"Uploads image",
			args{
				project:     "myproj",
				writeImages: true,
			},
			[]common.DataRecord{{
				Origin:         &common.OriginDescription{Name: "MyRac.rac"},
				RamsesHeader:   &ramses.Ramses{},
				RamsesTMHeader: &ramses.TMHeader{},
				SourceHeader:   &innosat.SourcePacketHeader{},
				TMHeader:       &innosat.TMHeader{},
				RID:            aez.CCD5,
				Data: &aez.CCDImage{
					PackData: &aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
					ImageFileName: "MyRac_5000000000_5.png",
				},
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			}},
			map[string]int{
				"myproj/MyRac_5000000000_5.png":  76,  // 8 + header
				"myproj/MyRac_5000000000_5.json": 857, // length of the json
			},
		},
		{
			"Doesn't upload image when told not to",
			args{
				project:     "myproj",
				writeImages: false,
			},
			[]common.DataRecord{{
				Origin:         &common.OriginDescription{Name: "MyRac.rac"},
				RamsesHeader:   &ramses.Ramses{},
				RamsesTMHeader: &ramses.TMHeader{},
				SourceHeader:   &innosat.SourcePacketHeader{},
				TMHeader:       &innosat.TMHeader{},
				Data: &aez.CCDImage{
					PackData: &aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
					ImageFileName: "MyRac_5000000000_5.png",
				},
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			}},
			map[string]int{},
		},
		{
			"Doesn't upload errors",
			args{
				project:     "myproj",
				writeImages: true,
			},
			[]common.DataRecord{{
				Origin:         &common.OriginDescription{Name: "MyRac.rac"},
				RamsesHeader:   &ramses.Ramses{},
				RamsesTMHeader: &ramses.TMHeader{},
				SourceHeader:   &innosat.SourcePacketHeader{},
				TMHeader:       &innosat.TMHeader{},
				Data: &aez.CCDImage{
					PackData: &aez.CCDImagePackData{
						EXPTS: 5,
						JPEGQ: aez.JPEGQUncompressed16bit,
						NCOL:  1,
						NROW:  2,
					},
					ImageFileName: "MyRac_5000000000_5.png",
				},
				Error:  errors.New("here be dragons"),
				Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
			}},
			map[string]int{},
		},
		{
			"Uploads everything",
			args{
				project:             "myproj",
				descriptionFileName: "info.json",
				descriptionFileBody: []byte("[42,42]"),
				writeImages:         true,
				writeTimeseries:     true,
			},
			[]common.DataRecord{
				{
					Origin:         &common.OriginDescription{Name: "MyRac.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					RID:            aez.CCD6,
					Data: &aez.CCDImage{
						PackData: &aez.CCDImagePackData{
							EXPTS: 5,
							JPEGQ: aez.JPEGQUncompressed16bit,
							NCOL:  1,
							NROW:  2,
						},
						ImageFileName: "MyRac_5000000000_6.png",
					},
					Buffer: make([]byte, 2*2*2), // 2x2 pixels, 2 bytes per pix
				},
				{
					Origin:         &common.OriginDescription{Name: "MyOtherRac.rac"},
					RamsesHeader:   &ramses.Ramses{},
					RamsesTMHeader: &ramses.TMHeader{},
					SourceHeader:   &innosat.SourcePacketHeader{},
					TMHeader:       &innosat.TMHeader{},
					Data:           &aez.HTR{},
				},
			},
			map[string]int{
				"myproj/ABOUT.json":              7,
				"myproj/MyRac_5000000000_6.png":  76,  // 8 + header
				"myproj/MyRac_5000000000_6.json": 857, // length of the json
				"myproj/CCD.csv":                 708, // length of the first three lines csv (specs, header, datarow)
				"myproj/HTR.csv":                 1005,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			var idxUp = 0

			var uploader = func(uploader *s3manager.Uploader, key string, bodyBuffer io.Reader) {
				key = strings.ReplaceAll(key, "\\", "/")
				buf, _ := ioutil.ReadAll(bodyBuffer)
				if idxUp >= len(tt.uploads) {
					t.Errorf(
						"Got unexpected upload #%v, key '%v', body %v",
						idxUp,
						key,
						buf,
					)
				} else {
					bodyLen, ok := tt.uploads[key]
					if !ok {
						t.Errorf("Upload %v: key = %v, key not wanted", idxUp, key)
					} else if bodyLen != len(buf) {
						t.Errorf("Upload %v/%v: len(buf) = %v, want %v ", idxUp, key, len(buf), bodyLen)
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

			for _, record := range tt.records {
				callback(record)
			}
			teardown()

			if idxUp < len(tt.uploads) {
				t.Errorf("Recorded %v uploads, want %v", idxUp, len(tt.uploads))
			}
		})
	}
}

func Test_csvAWSWriterFactoryCreator(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("localhost")}))
	upload := s3manager.NewUploader(sess)
	uploads := make(map[string]int)

	var uploader = func(uploader *s3manager.Uploader, key string, bodyBuffer io.Reader) {
		buf, _ := ioutil.ReadAll(bodyBuffer)
		uploads[key] = len(buf)
	}
	factory := csvAWSWriterFactoryCreator(uploader, upload, "myproject")
	writer, err := factory(&common.DataRecord{Data: &aez.HTR{}}, timeseries.HTR)
	if err != nil {
		t.Errorf("csvAWSWriterFactoryCreator()'s factory returned error %v", err)
	}
	writer.Close()
	want := "myproject/HTR.csv"
	_, ok := uploads[want]
	if !ok {
		t.Errorf("csvAWSWriterFactoryCreator()'s factory produced uploads %v, want key %v", uploads, want)
	}
}
