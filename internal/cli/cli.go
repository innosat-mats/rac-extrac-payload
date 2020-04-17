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
		err = ramsesData.Read(f)
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

		reader := bytes.NewReader(ramsesPayload)
		if ramsesData.SecureTrans() {
			ramsesSecure := ramses.Secure{}
			err = ramsesSecure.Read(reader)
			if err != nil {
				log.Fatal(err)
			}
			startPayload = binary.Size(ramsesSecure)
		}

		// Source packet header
		header := innosat.SourcePacketHeader{}
		err = header.Read(reader)
		if err != nil {
			log.Fatal(err)
		}
		if crc16.ChecksumCCITTFalse(ramsesPayload[startPayload:ramsesData.Length-2]) != binary.BigEndian.Uint16(ramsesPayload[ramsesData.Length-2:]) {
			log.Fatal("checksum bad")
		}

		if header.Type() != innosat.TC {
			// Telecommands not supported yet
			continue
		}
		// Data Field Header
		dataHeader := innosat.TMDataFieldHeader{}
		err = dataHeader.Read(reader)
		if err != nil {
			log.Fatal(err)
		}

		if header.Type() == innosat.TM && header.IsMainApplication() && dataHeader.IsHousekeeping() {
			var sid aez.SID
			binary.Read(reader, binary.BigEndian, &sid)
			if sid == aez.SIDSTAT {
				stat := aez.STAT{}
				err = stat.Read(reader)
				if err != nil {
					log.Fatal("stat", err)
				}
			}
			if sid == aez.SIDHTR {
				htr := aez.HTR{}
				err = htr.Read(reader)
				if err != nil {
					log.Fatal("htr", err)
				}
			}
			if sid == aez.SIDPWR {
				pwr := aez.PWR{}
				err = pwr.Read(reader)
				if err != nil {
					log.Fatal("pwr", err)
				}
			}
			if sid == aez.SIDCPRUA {
				cprua := aez.CPRU{}
				err = cprua.Read(reader)
				if err != nil {
					log.Fatal("cprua", err)
				}
			}
			if sid == aez.SIDCPRUB {
				cprub := aez.CPRU{}
				err = cprub.Read(reader)
				if err != nil {
					log.Fatal("cprub", err)
				}
			}

		}

	}

}
