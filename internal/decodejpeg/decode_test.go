package decodejpeg

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image"
	"io/ioutil"
	"log"

	pnm "github.com/jbuchbinder/gopnm"
)

func unWindowTrunc(pixValue uint16) uint16 {
	return pixValue << 4
}

func Example() {
	fileContents, err := ioutil.ReadFile("testdata/3166_4052_5.jpg")
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
	fmt.Printf("Jpeg Data md5sum:      %x\n", md5.Sum(img.Pix))

	ReferenceData, err := ioutil.ReadFile("testdata/3166_4052_5.pnm")
	if err != nil {
		log.Fatalln(err)
	}
	rreader := bytes.NewReader(ReferenceData)
	refImg, err := pnm.Decode(rreader)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Reference Data md5sum: %x\n", md5.Sum(refImg.(*image.Gray16).Pix))
	//Output:
	//Jpeg Data md5sum:      d43342ed75032b8f7ae380bb1ff78459
	//Reference Data md5sum: d43342ed75032b8f7ae380bb1ff78459
}
