package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	pus "github.com/innosat-mats/rac-extract-payload/internal/innosatPUS"
)

func main() {
	fmt.Println("hej")
	var err error
	var n int
	var messageType uint
	var buf []byte
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < 20; i++ {
		// Read RAC PDU part
		pduData := make([]byte, 32)
		n, err = f.Read(pduData)
		if err != nil {
			log.Fatal(err, n)
		}

		// Source packet header

		header := pus.SourcePacketHeader{}
		err = header.Read(f)
		if err != nil {
			log.Panic(err)
		}

		// packetID
		fmt.Printf("%03b ", header.Version())
		fmt.Printf("%b ", header.Type())
		fmt.Printf("%b ", header.HeaderType())
		fmt.Printf("%5d ", header.APID())

		//sequence
		fmt.Printf("%02b ", header.GroupingFlags())
		fmt.Printf("%5d ", header.SequenceCount())

		//size
		fmt.Printf("%7d ", header.PacketLength)

		buf = make([]byte, header.PacketLength+1)
		n, err = f.Read(buf)
		if err != nil {
			log.Fatal(err, n)
		}

		// Data Field Header
		dataHeader := pus.TMDataFieldHeader{}
		dhsize := binary.Size(dataHeader)
		dhr := bytes.NewReader(buf[0:])
		err = dataHeader.Read(dhr)
		if err != nil {
			log.Fatal(err)
		}

		// Data field  header
		fmt.Printf("%04b ", dataHeader.PUSVersion())

		fmt.Printf("%4d ", dataHeader.ServiceType)
		fmt.Printf("%4d ", dataHeader.ServiceSubType)
		if messageType == 0 {
			fmt.Printf("%d ", dataHeader.Time())
		}

		// if crc16.ChecksumCCITTFalse(append(header, data[:dataPackageLength-2]...)) == binary.BigEndian.Uint16(data[dataPackageLength-2:dataPackageLength]) {
		// 	fmt.Print("✔️ ")
		// } else {
		// 	fmt.Print("❌ ")
		// }
		if dataHeader.ServiceType == 3 && dataHeader.ServiceSubType == 25 {
			sid := binary.BigEndian.Uint16(buf[dhsize:])
			r := bytes.NewReader(buf[dhsize+2:])
			if sid == 1 {
				stat := pus.STAT{}
				err = stat.Read(r)
				if err != nil {
					log.Fatal("stat", err)
				}
				fmt.Print(" STAT: ", stat)
			}
			if sid == 10 {
				htr := pus.HTR{}
				err = htr.Read(r)
				// if err != nil {
				// 	log.Fatal("htr", err)
				// }
				fmt.Print(" HTR:  ", htr)
			}
			if sid == 20 {
				pwr := pus.PWR{}
				err = pwr.Read(r)
				if err != nil {
					log.Fatal("pwr", err)
				}
				fmt.Print(" PWR:  ", pwr)
			}
			if sid == 30 {
				cprua := pus.CPRU{}
				err = cprua.Read(r)
				if err != nil {
					log.Fatal("cprua", err)
				}
				fmt.Print(" CPRUA:", cprua)
			}
			if sid == 31 {
				cprub := pus.CPRU{}
				err = cprub.Read(r)
				if err != nil {
					log.Fatal("cprub", err)
				}
				fmt.Print(" CPRUB:", cprub)
			}

		}

		fmt.Println()
	}

}
