package main

import (
	"fmt"
	"log"
	"os"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	common.ExtractData(f, fmt.Println)
}
