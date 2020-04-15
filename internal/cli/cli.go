package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/howeyc/crc16"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func main() {
	var err error
	var startPayload int
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for {
		// Read ramses part
		ramsesData := ramses.Ramses{}
		_, err = ramsesData.Read(f)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if !ramsesData.Valid() {
			log.Fatal("Not a valid RAC-file")
		}

		// Get the payload inside the ramses packate
		ramsesPayload := make([]byte, ramsesData.Length)
		_, err = f.Read(ramsesPayload)
		if err != nil {
			log.Fatal(err)
		}

		r := bytes.NewReader(ramsesPayload)
		if ramsesData.SecureTrans() {
			ramsesSecure := ramses.Secure{}
			startPayload, err = ramsesSecure.Read(r)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Source packet header

		header := innosat.SourcePacketHeader{}
		err = header.Read(r)
		if err != nil {
			log.Fatal(err)
		}

		// Data Field Header
		dataHeader := innosat.TMDataFieldHeader{}
		err = dataHeader.Read(r)
		if err != nil {
			log.Fatal(err)
		}

		if crc16.ChecksumCCITTFalse(ramsesPayload[startPayload:ramsesData.Length-2]) != binary.BigEndian.Uint16(ramsesPayload[ramsesData.Length-2:]) {
			log.Fatal("checksum bad")
		}
		if header.Type() == 1 && header.APID() == 100 && dataHeader.ServiceType == 3 && dataHeader.ServiceSubType == 25 {
			var sid uint16
			binary.Read(r, binary.BigEndian, &sid)
			if sid == 1 {
				stat := aez.STAT{}
				err = stat.Read(r)
				if err != nil {
					log.Fatal("stat", err)
				}
			}
			if sid == 10 {
				htr := aez.HTR{}
				err = htr.Read(r)
				if err != nil {
					log.Fatal("htr", err)
				}
			}
			if sid == 20 {
				pwr := aez.PWR{}
				err = pwr.Read(r)
				if err != nil {
					log.Fatal("pwr", err)
				}
			}
			if sid == 30 {
				cprua := aez.CPRU{}
				err = cprua.Read(r)
				if err != nil {
					log.Fatal("cprua", err)
				}
			}
			if sid == 31 {
				cprub := aez.CPRU{}
				err = cprub.Read(r)
				if err != nil {
					log.Fatal("cprub", err)
				}
			}

		}

	}

}
