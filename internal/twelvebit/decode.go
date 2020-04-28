package main

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: ${SRCDIR}/lib/libjpeg.a
// #include <stdio.h>
// #include <stdlib.h>
// int read_JPEG_file (char*);
import "C"
import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"unsafe"
)

var row = 20
var col = 20

var buffer = make([]byte, 0, 20*20)
var imageData = bytes.NewBuffer(buffer)

func main() {
	cs := C.CString(os.Args[1])
	C.read_JPEG_file(cs)
	C.free(unsafe.Pointer(cs))
	img := image.NewGray(image.Rect(0, 0, row, col))
	img.Pix = imageData.Bytes()
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
}

//export callback
func callback(buf unsafe.Pointer, length C.int) {
	scanline := C.GoBytes(buf, length)
	imageData.Write(scanline)
}
