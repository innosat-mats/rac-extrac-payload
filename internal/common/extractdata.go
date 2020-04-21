package common

import (
	"io"
	"log"
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// ExtractData reads Ramses data packages and extract the instrument data.
func ExtractData(byteStream io.Reader, extract func(aez ...interface{}) (int, error)) {
	var waitGroup sync.WaitGroup
	ramsesPackages := make(chan ramses.Package)
	ramsesErrors := make(chan error)
	innosatPackages := make(chan innosat.SourcePackage)
	dataPackages := make(chan aez.PackageType)
	go aez.DecodePackages(innosatPackages, dataPackages)
	go ramses.Packages(byteStream, ramsesPackages, ramsesErrors)

	go func() {
		waitGroup.Add(1)
		for data := range dataPackages {
			extract(data)
		}
		waitGroup.Done()
	}()

	for ramsesPackage := range ramsesPackages {
		select {
		case err := <-ramsesErrors:
			if err != nil {
				log.Output(log.Llongfile, err.Error())
				continue
			}
		default:
			innosatPackage := innosat.DecodeSource(ramsesPackage)
			innosatPackages <- innosatPackage
		}
	}
	close(innosatPackages)
	waitGroup.Wait()
}
