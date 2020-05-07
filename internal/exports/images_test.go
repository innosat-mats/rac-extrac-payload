package exports

import (
	"image"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
)

func Test_getGrayscaleImageName(t *testing.T) {
	type args struct {
		dir         string
		imgPackData aez.CCDImagePackData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Uses nanos in filename",
			args{"test", aez.CCDImagePackData{EXPTS: 5}},
			filepath.Join("test", "5000000000.png"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGrayscaleImageName(tt.args.dir, tt.args.imgPackData); got != tt.want {
				t.Errorf("getGrayscaleImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestImage() []byte {
	fileContents, err := ioutil.ReadFile(filepath.Join("testdata", "3166_4052_5.jpg"))
	if err != nil {
		log.Fatalln(err)
	}
	return fileContents
}

func Test_getImageData(t *testing.T) {
	// NON-STANDARD test because don't really want to check full image
	type args struct {
		buf         []byte
		packData    aez.CCDImagePackData
		outFileName string
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
		want       []uint16
	}{
		{
			"Processes uncompressed directly as pixels",
			args{
				buf:         []byte{0xff, 0x00},
				packData:    aez.CCDImagePackData{JPEGQ: aez.JPEGQUncompressed16bit},
				outFileName: "myfile.png",
			},
			1,
			[]uint16{255},
		},
		{
			"Processes compressed jpeg 12bit buffer into pixels",
			args{
				buf:         getTestImage(),
				packData:    aez.CCDImagePackData{JPEGQ: 95},
				outFileName: "myfile.png",
			},
			250 * 501,
			[]uint16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getImageData(tt.args.buf, tt.args.packData, tt.args.outFileName)
			if len(got) != tt.wantLength {
				t.Errorf("getImageData() returned %v pixels, want %v", len(got), tt.wantLength)
			}
			if len(tt.want) > 0 && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImageData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestImagePixels(packData aez.CCDImagePackData) []uint16 {
	buf := getTestImage()
	return getImageData(buf, packData, "test.png")
}

func Test_getGrayscaleImage(t *testing.T) {
	type args struct {
		pixels   []uint16
		width    int
		height   int
		shift    int
		filename string
	}
	tests := []struct {
		name        string
		args        args
		wantShape   image.Rectangle
		samplePixel image.Point
		wantValue   uint16
	}{
		{
			"Returns expected image unshifted",
			args{
				pixels:   getTestImagePixels(aez.CCDImagePackData{JPEGQ: 95}),
				width:    501,
				height:   250,
				shift:    0,
				filename: "myfile.png",
			},
			image.Rectangle{image.Point{0, 0}, image.Point{501, 250}},
			image.Point{42, 1},
			317,
		},
		{
			"Returns expected image shifted",
			args{
				pixels:   getTestImagePixels(aez.CCDImagePackData{JPEGQ: 95}),
				width:    501,
				height:   250,
				shift:    1,
				filename: "myfile.png",
			},
			image.Rectangle{image.Point{0, 0}, image.Point{501, 250}},
			image.Point{42, 1},
			317 * 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getGrayscaleImage(
				tt.args.pixels,
				tt.args.width,
				tt.args.height,
				tt.args.shift,
				tt.args.filename,
			)
			if !reflect.DeepEqual(got.Bounds(), tt.wantShape) {
				t.Errorf(
					"getGrayscaleImage() returns bounds %v, want %v", got.Bounds(), tt.wantShape,
				)
			}
			if pix := got.Gray16At(tt.samplePixel.X, tt.samplePixel.Y).Y; pix != tt.wantValue {
				t.Errorf(
					"getGrayscaleImage() returns image with pixel %v value %v, want %v",
					tt.samplePixel,
					pix,
					tt.wantValue,
				)

			}
		})
	}
}
