package common

import (
	"fmt"
	"io"
	"log"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// ExtractData ...
func ExtractData(byteStream io.Reader, extract func(aez ...interface{}) (int, error)) {

	ramsesPackages := make(chan ramses.Package)
	ramsesErrors := make(chan error)
	innosatPackages := make(chan innosat.SourcePackage)
	dataPackages := make(chan aez.PackageType, 1)

	go aez.DecodePackages(innosatPackages, dataPackages)
	go ramses.Packages(byteStream, ramsesPackages, ramsesErrors)
	for ramsesPackage := range ramsesPackages {
		select {
		case err := <-ramsesErrors:
			if err != nil {
				log.Fatalln(err)
			}
		case data, more := <-dataPackages:
			if more {
				extract(data)
			}
		default:
			innosatPackage := innosat.DecodeSource(ramsesPackage)
			innosatPackages <- innosatPackage
		}
	}
	fmt.Println()
}
