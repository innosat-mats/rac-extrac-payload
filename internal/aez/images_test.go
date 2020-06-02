package aez

import (
	"image"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getGrayscaleImageName(t *testing.T) {
	tests := []struct {
		name       string
		project    string
		originName string
		expts      uint32
		rid        RID
		want       string
	}{
		{
			"Adds project",
			"test",
			"my/path/MyFile.rac",
			5,
			CCD1,
			filepath.Join("test", "MyFile_5000000000_1.png"),
		},
		{
			"Omits project if empty",
			"",
			"my/path/MyFile.rac",
			5,
			CCD7,
			"MyFile_5000000000_7.png",
		},
	}
	for _, tt := range tests {
		got := getGrayscaleImageName(
			tt.project,
			tt.originName,
			&CCDImagePackData{EXPTS: tt.expts},
			tt.rid,
		)
		if got != tt.want {
			t.Errorf("getGrayscaleImageName() = %v, want %v", got, tt.want)
		}
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
		packData    CCDImagePackData
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
				packData:    CCDImagePackData{JPEGQ: JPEGQUncompressed16bit},
				outFileName: "myfile.png",
			},
			1,
			[]uint16{255},
		},
		{
			"Processes compressed jpeg 12bit buffer into pixels",
			args{
				buf:         getTestImage(),
				packData:    CCDImagePackData{JPEGQ: 95},
				outFileName: "myfile.png",
			},
			250 * 501,
			[]uint16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getImageData(tt.args.buf, &tt.args.packData, tt.args.outFileName)
			if len(got) != tt.wantLength {
				t.Errorf("getImageData() returned %v pixels, want %v", len(got), tt.wantLength)
			}
			if len(tt.want) > 0 && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImageData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestImagePixels(packData CCDImagePackData) []uint16 {
	buf := getTestImage()
	return getImageData(buf, &packData, "test.png")
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
				pixels:   getTestImagePixels(CCDImagePackData{JPEGQ: 95}),
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
				pixels:   getTestImagePixels(CCDImagePackData{JPEGQ: 95}),
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
