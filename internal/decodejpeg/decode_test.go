package decodejpeg

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"image"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	pnm "github.com/jbuchbinder/gopnm"
)

func TestJpegImage(t *testing.T) {
	fileContents, err := ioutil.ReadFile("testdata/lava.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	raw, width, height, err := JpegImageData(fileContents)
	bufArray := []byte{}
	buf := bytes.NewBuffer(bufArray)
	err = binary.Write(buf, binary.BigEndian, raw)
	if err != nil {
		log.Fatalln(err)
	}
	img := image.NewGray16(image.Rect(0, 0, width, height))
	img.Pix = buf.Bytes()

	ReferenceData, err := ioutil.ReadFile("testdata/lava.pnm")
	if err != nil {
		log.Fatalln(err)
	}
	rreader := bytes.NewReader(ReferenceData)
	refImg, err := pnm.Decode(rreader)
	if err != nil {
		log.Fatalln(err)
	}
	if md5.Sum(img.Pix) != md5.Sum(refImg.(*image.Gray16).Pix) {
		t.Errorf("JpegImageData() integration test: Image data is the same")
	}

}

func getData(base64String string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return []byte{}
	}
	return decoded
}

func genData(value uint16, size int) []uint16 {
	slice := make([]uint16, size)
	for i := range slice {
		slice[i] = value

	}
	return slice
}

func TestJpegImageData(t *testing.T) {
	type args struct {
		jpegData []byte
	}
	tests := []struct {
		name        string
		args        args
		wantRawData []uint16
		wantHeight  int
		wantWidth   int
		wantErr     bool
	}{
		{
			"all white",
			args{getData(white20x20jpg)},
			genData(1<<12-1, 400),
			20,
			20,
			false,
		},
		{
			"all white large",
			args{getData(white500x500jpg)},
			genData(1<<12-1, 500*500),
			500,
			500,
			false,
		},
		{
			"all black",
			args{getData(black20x20jpg)},
			genData(0, 400),
			20,
			20,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRawData, gotHeight, gotWidth, err := JpegImageData(tt.args.jpegData)
			if (err != nil) != tt.wantErr {
				t.Errorf("JpegImageData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRawData, tt.wantRawData) {
				t.Errorf("JpegImageData() gotRawData = %v, want %v", gotRawData, tt.wantRawData)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("JpegImageData() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("JpegImageData() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
		})
	}
}
func BenchmarkJpegImageData(b *testing.B) {
	whiteData := getData(white500x500jpg)
	refData := genData(1<<12-1, 500*500)

	for i := 0; i < b.N; i++ {
		gotData, _, _, err := JpegImageData(whiteData)
		if err != nil {
			b.Errorf("JpegImageData() failed with: %s", err)
		}
		if !reflect.DeepEqual(refData, gotData) || err != nil {
			b.Errorf("JpegImageData() refData = %v, want %v", refData, gotData)
		}

	}
}

const white20x20jpg = `/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0a
HBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/yQALDAAUABQBAREA/8wABgAQEAX/2gAI
AQEAAD8A0u28jP/Z`

const black20x20jpg = `/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0a
HBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/yQALDAAUABQBAREA/8wABgAQEAX/2gAI
AQEAAD8A/wD8HeD/2Q==`

const white500x500jpg = `/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0a
HBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/yQALDAH0AfQBAREA/8wABgAQEAX/2gAI
AQEAAD8A0u28i2Cg/9k=`
