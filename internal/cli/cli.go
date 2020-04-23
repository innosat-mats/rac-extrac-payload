package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func callback(pkg common.DataRecord) {
	fmt.Println(pkg)
}

func main() {
	batch := make([]common.StreamBatch, len(os.Args)-1)
	for n, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}

		batch[n] = common.StreamBatch{
			Buf: f,
			Origin: common.OriginDescription{
				Name: filename,
				Date: time.Now(),
			},
		}

	}
	common.ExtractData(callback, batch...)
}
