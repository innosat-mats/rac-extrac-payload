package main

// #cgo CFLAGS: -I${SRCDIR}/jpeglib-8bit/include
// #cgo LDFLAGS: ${SRCDIR}/jpeglib-8bit/lib/libjpeg.a
// #include <stdlib.h>
// char* read_JPEG_file (char*, size_t);
import "C"
import (
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"unsafe"
)

func main() {
	width := 20
	height := 20
	imageFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	rawData := C.CString(string(imageFile))
	defer C.free(unsafe.Pointer(rawData))

	pixelData := C.read_JPEG_file(rawData, C.size_t(len(imageFile)))
	defer C.free(unsafe.Pointer(pixelData))

	img := image.NewGray(image.Rect(0, 0, width, height))
	img.Pix = C.GoBytes(unsafe.Pointer(pixelData), C.int(width)*C.int(height))

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
