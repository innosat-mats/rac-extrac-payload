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
		0x90, 0xeb, 0x48, 0x00, 0x79,
		0xd8, 0x00, 0x00, 0xfb, 0xad, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x02, 0x00,
		0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00, 0x08, 0x64, 0xc8, 0x98, 0x00,
		0x31, 0x10, 0x03, 0x19, 0x00, 0x00, 0x12, 0x19, 0xe3, 0x39, 0x00, 0x01, 0x7f, 0x04, 0x02, 0x82,
		0x04, 0x02, 0x02, 0x06, 0x01, 0x19, 0x12, 0x00, 0x00, 0x0c, 0xe3, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x30, 0xfb,
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
	// 2020-03-22 22:17:40.603 +0000 UTC
	// true
	// {32516 2 33284 2 2 6 1 420610048 3299 2 0 0 16777216 1090519040 0 0}
}
