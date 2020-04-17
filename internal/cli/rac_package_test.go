package main_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/howeyc/crc16"
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// Example ...
func Example() {
	data := []byte{
		0x90, 0xeb, 0x45, 0x00, 0x79, 0xd8, 0x00, 0x00, 0x54, 0x4b, 0x3f, 0x03, 0x57, 0x1b, 0x00, 0x00,
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		0x08, 0x64, 0xe1, 0x19, 0x00, 0x2e, 0x10, 0x03, 0x19, 0x00, 0x00, 0x26, 0x94, 0x99, 0x1c, 0x00,
		0x01, 0x7f, 0x04, 0x02, 0x82, 0x04, 0x02, 0x94, 0x26, 0x00, 0x00, 0x12, 0x99, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x6c, 0x00, 0x00, 0x00, 0x19, 0x00,
		0x00, 0x00, 0x00, 0x6a, 0x65,
	}

	var err error
	byteStream := bytes.NewReader(data)
	for {
		// Read ramses part
		ramsesData := ramses.Ramses{}
		err = ramsesData.Read(byteStream)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fmt.Println(ramsesData.Created())

		if !ramsesData.Valid() {
			log.Fatal("Not a valid RAC-file")
		}

		// Get the payload inside the ramses packate
		ramsesPayload := make([]byte, ramsesData.Length)
		_, err = byteStream.Read(ramsesPayload)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(ramsesPayload)
		var startPayload int
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

		fmt.Println(header.IsMainApplication())

		dataHeader := innosat.TMDataFieldHeader{}
		err = dataHeader.Read(reader)
		if err != nil {
			log.Fatal(err)
		}

		if header.IsMainApplication() && dataHeader.IsHousekeeping() {
			var sid uint16
			binary.Read(reader, binary.BigEndian, &sid)
			if sid == 1 {
				stat := aez.STAT{}
				err = stat.Read(reader)
				if err != nil {
					log.Fatal("stat", err)
				}
				fmt.Println(stat)
			}
		}
	}
	// Output:
	// 2019-03-01 15:07:59.7 +0000 UTC
	// true
	// {32516 2 33284 2 2485518336 4761 2 0 0 50331648 1811939328 419430400}
}
