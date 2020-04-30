package djpeg

import (
	"fmt"
	"io/ioutil"
	"log"
)

func unWindowTrunc(pixValue uint16) uint16 {
	return pixValue << 4
}

func Example() {
	fileContents, err := ioutil.ReadFile("3166_4052_5.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	img := DecodeJpegContents(fileContents, 501, 200, unWindowTrunc)
	fmt.Println(img.Pix[:10])
	// Output: [14 80 15 16 16 0 17 64 18 48]
}
