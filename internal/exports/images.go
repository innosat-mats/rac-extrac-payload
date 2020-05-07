package exports

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/decodejpeg"
)

func getGrayscaleImage(
	pixels []uint16, width int, height int, shift int, filename string,
) image.Image {
	nPixels := len(pixels)
	/*
		if nPixels == 1048572 {
			log.Printf(
				"changing widht %v => 1026 and height %v => 1022",
				width,
				height,
			)
			width = 1026
			height = 1022
		}
	*/
	img := image.NewGray16(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	if nPixels != width*height {
		log.Printf(
			"%v Found %v pixels, but dimension %v x %v says it should be %v",
			filename,
			nPixels,
			width,
			height,
			width*height,
		)
	}
	shifted := make([]uint16, nPixels)
	for idx := 0; idx < nPixels; idx++ {
		shifted[idx] = pixels[idx] << shift
	}
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, shifted)
	if err != nil {
		log.Print("Could not write image data")
		return img
	}
	img.Pix = buf.Bytes()

	/*

		nPixels := len(pixels)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				idx := x*height + y
				if idx < nPixels {
					pixel := pixels[idx] << shift
					img.Set(x, y, color.Gray16{pixel})
				}
			}
		}
	*/
	return img
}

func getGrayscaleImageName(dir string, imgPackData aez.CCDImagePackData) string {
	return filepath.Join(
		dir,
		fmt.Sprintf("%v.png", imgPackData.Nanoseconds()),
	)
}

func getImageData(
	buf []byte,
	packData aez.CCDImagePackData,
	outFileName string,
) []uint16 {
	var imgData []uint16
	var err error
	if packData.JPEGQ != aez.JPEGQUncompressed16bit {
		log.Println("Compressed image")
		var height int
		var width int
		imgData, height, width, err = decodejpeg.JpegImageData(buf)
		if err != nil {
			log.Print(err)
			return imgData
		}
		if uint16(height) != packData.NROW || uint16(width) != packData.NCOL+1 {
			log.Printf(
				"CCDImage %v has either width %v != %v and/or height %v != %v\n",
				outFileName,
				height,
				packData.NROW,
				width,
				packData.NCOL,
			)
		}
	} else {
		log.Println("Raw image")
		reader := bytes.NewReader(buf)
		imgData = make([]uint16, reader.Len()/2)
		binary.Read(reader, binary.LittleEndian, &imgData)
	}
	return imgData
}
