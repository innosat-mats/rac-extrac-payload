package djpeg

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: ${SRCDIR}/lib/libjpeg.a
// #include <stdlib.h>
// char* read_JPEG_file (char*, size_t);
import "C"
import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"log"
	"unsafe"
)

//DecodeJpegContents converts an image encoded in 12-bit jpeg to an image of Gray16 type
func DecodeJpegContents(jpegData []byte, width int, height int, untrunc func(pix uint16) uint16) *image.Gray16 {

	rawData := C.CString(string(jpegData))
	defer C.free(unsafe.Pointer(rawData))

	pixelData, err := C.read_JPEG_file(rawData, C.size_t(len(jpegData)))
	if err != nil {
		log.Fatalln(err)
	}
	defer C.free(unsafe.Pointer(pixelData))

	pixelString := C.GoStringN(pixelData, C.int(width)*C.int(height)*2)

	buffer := bytes.NewBufferString(pixelString)

	img := image.NewGray16(image.Rect(0, 0, width, height))
	imageBuffer := bytes.NewBuffer(img.Pix)
	imageBuffer.Reset()

	for {
		var pixel uint16
		var err error
		err = binary.Read(buffer, binary.LittleEndian, &pixel)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		pixel = untrunc(pixel)
		err = binary.Write(imageBuffer, binary.BigEndian, &pixel)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return img
}
