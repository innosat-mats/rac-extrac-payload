package decodejpeg

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: ${SRCDIR}/lib/libjpeg.a
// #include <stdlib.h>
// #include "decode.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//JpegImageData converts a grayscale image encoded in 12-bit jpeg to raw data
func JpegImageData(jpegData []byte) ([]uint16, int, int, error) {

	rawData := C.CString(string(jpegData))
	defer C.free(unsafe.Pointer(rawData))
	pixelData, err := C.read_JPEG_file(rawData, C.size_t(len(jpegData)))
	if err != nil {
		return []uint16{}, 0, 0, err
	}
	defer C.free(unsafe.Pointer(pixelData.pix))

	pixelString := C.GoStringN(pixelData.pix, C.int(pixelData.width*pixelData.height*C.sizeof_short))

	buffer := bytes.NewBufferString(pixelString)

	outputData := make([]uint16, pixelData.width*pixelData.height)

	err = binary.Read(buffer, binary.LittleEndian, outputData)
	if err != nil {
		return []uint16{}, 0, 0, err
	}
	return outputData, int(pixelData.height), int(pixelData.width), nil
}
