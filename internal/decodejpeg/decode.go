package decodejpeg

// #cgo windows CFLAGS: -I${SRCDIR}/../../third-party/windows/include
// #cgo windows LDFLAGS: ${SRCDIR}/../../third-party/windows/lib/libjpeg.a
// #cgo linux CFLAGS: -I${SRCDIR}/../../third-party/linux/include
// #cgo linux LDFLAGS: ${SRCDIR}/../../third-party/linux/lib/libjpeg.a
// #cgo darwin CFLAGS: -I${SRCDIR}/../../third-party/darwin/include
// #cgo darwin LDFLAGS: -L/${SRCDIR}/../../third-party/darwin/lib -Wl,-rpath,$ORIGIN/../lib -ljpeg
// #include <stdlib.h>
// #include "decode.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//JpegImageData converts a grayscale image encoded in 12-bit jpeg to raw data
func JpegImageData(jpegData []byte) (rawData []uint16, height int, width int, err error) {

	jpegChar := C.CString(string(jpegData))
	defer C.free(unsafe.Pointer(jpegChar))

	imageData, err := C.read_JPEG_file(jpegChar, C.size_t(len(jpegData)))
	if err != nil {
		return rawData, height, width, err
	}
	defer C.free(unsafe.Pointer(imageData.pix))

	pixelString := C.GoStringN(imageData.pix, C.int(imageData.width*imageData.height*C.sizeof_short))
	height = int(imageData.height)
	width = int(imageData.width)
	buffer := bytes.NewBufferString(pixelString)
	rawData = make([]uint16, width*height)
	err = binary.Read(buffer, binary.LittleEndian, rawData)
	if err != nil {
		return rawData, height, width, err
	}
	return rawData, height, width, nil
}
