package aez

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

//PackageType ...
type PackageType interface {
}

// DecodePackages ...
func DecodePackages(source chan innosat.SourcePackage, target chan PackageType) {
	defer close(target)
	for sourcePackage := range source {

		var err error
		var data PackageType
		reader := bytes.NewReader(sourcePackage.Application)
		if sourcePackage.Header.GroupingFlags() == innosat.SPStandalone {

			if sourcePackage.Header.Type() == innosat.TM && sourcePackage.Header.IsMainApplication() {
				var sid SID
				binary.Read(reader, binary.BigEndian, &sid)
				if sid == SIDSTAT {
					stat := STAT{}
					err = stat.Read(reader)
					if err != nil {
						log.Fatal("stat", err)
					}
					data = stat
				}
				if sid == SIDHTR {
					htr := HTR{}
					err = htr.Read(reader)
					if err != nil {
						log.Fatal("htr", err)
					}
					data = htr
				}
				if sid == SIDPWR {
					pwr := PWR{}
					err = pwr.Read(reader)
					if err != nil {
						log.Fatal("pwr", err)
					}
					data = pwr
				}
				if sid == SIDCPRUA {
					cpru := CPRU{}
					err = cpru.Read(reader)
					if err != nil {
						log.Fatal("cprua", err)
					}
					data = cpru
				}
				if sid == SIDCPRUB {
					cpru := CPRU{}
					err = cpru.Read(reader)
					if err != nil {
						log.Fatal("cprub", err)
					}
					data = cpru
				}
				target <- data
			}
		}

	}

}
