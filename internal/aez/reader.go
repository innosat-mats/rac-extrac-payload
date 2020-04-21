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

// Packages ...
func Packages(source <-chan innosat.SourcePackage, target chan<- PackageType) {
	defer close(target)
	for sourcePackage := range source {
		var err error
		var data PackageType
		reader := bytes.NewReader(sourcePackage.Application)
		if sourcePackage.Header.GroupingFlags() == innosat.SPStandalone {
			if sourcePackage.Header.Type() == innosat.TM && sourcePackage.Header.IsMainApplication() {
				var sid SID
				binary.Read(reader, binary.BigEndian, &sid)
				switch sid {
				case SIDSTAT:
					stat := STAT{}
					err = stat.Read(reader)
					if err != nil {
						log.Output(log.Llongfile, err.Error())
					}
					data = stat
				case SIDHTR:
					htr := HTR{}
					err = htr.Read(reader)
					if err != nil {
						log.Output(log.Llongfile, err.Error())
					}
					data = htr
				case SIDPWR:
					pwr := PWR{}
					err = pwr.Read(reader)
					if err != nil {
						log.Output(log.Llongfile, err.Error())
					}
					data = pwr
				case SIDCPRUA:
					cpru := CPRU{}
					err = cpru.Read(reader)
					if err != nil {
						log.Output(log.Llongfile, err.Error())
					}
					data = cpru
				case SIDCPRUB:
					cpru := CPRU{}
					err = cpru.Read(reader)
					if err != nil {
						log.Output(log.Llongfile, err.Error())
					}
					data = cpru
				}
				target <- data
			}
		}
	}
}
