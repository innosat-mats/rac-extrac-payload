package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func main() {
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

		header := SourcePacketHeader{}
		err = header.read(f)
		if err != nil {
			log.Panic(err)
		}

		// packetID
		fmt.Printf("%03b ", header.getVersion())
		fmt.Printf("%b ", header.getType())
		fmt.Printf("%b ", header.getHeaderType())
		fmt.Printf("%5d ", header.getAPID())

		//sequence
		fmt.Printf("%02b ", header.getGroupingFlags())
		fmt.Printf("%5d ", header.getSequenceCount())

		//size
		fmt.Printf("%7d ", header.PacketLength)

		buf = make([]byte, header.PacketLength+1)
		n, err = f.Read(buf)
		if err != nil {
			log.Fatal(err, n)
		}

		// Data Field Header
		dataHeader := TMDataFieldHeader{}
		dhsize := binary.Size(dataHeader)
		dhr := bytes.NewReader(buf[0:])
		err = dataHeader.read(dhr)
		if err != nil {
			log.Fatal(err)
		}

		// Data field  header
		fmt.Printf("%04b ", dataHeader.getPUS())

		fmt.Printf("%4d ", dataHeader.getServiceType())
		fmt.Printf("%4d ", dataHeader.getServiceSubType())
		if messageType == 0 {
			fmt.Printf("%d ", dataHeader.getTime())
		}

		// if crc16.ChecksumCCITTFalse(append(header, data[:dataPackageLength-2]...)) == binary.BigEndian.Uint16(data[dataPackageLength-2:dataPackageLength]) {
		// 	fmt.Print("✔️ ")
		// } else {
		// 	fmt.Print("❌ ")
		// }
		if dataHeader.getServiceType() == 3 && dataHeader.getServiceSubType() == 25 {
			sid := binary.BigEndian.Uint16(buf[dhsize:])
			r := bytes.NewReader(buf[dhsize+2:])
			if sid == 1 {
				stat := STAT{}
				err = stat.read(r)
				if err != nil {
					log.Fatal("stat", err)
				}
				fmt.Print(" STAT: ", stat)
			}
			if sid == 10 {
				htr := HTR{}
				err = htr.read(r)
				// if err != nil {
				// 	log.Fatal("htr", err)
				// }
				fmt.Print(" HTR:  ", htr)
			}
			if sid == 20 {
				pwr := PWR{}
				err = pwr.read(r)
				if err != nil {
					log.Fatal("pwr", err)
				}
				fmt.Print(" PWR:  ", pwr)
			}
			if sid == 30 {
				cprua := CPRU{}
				err = cprua.read(r)
				if err != nil {
					log.Fatal("cprua", err)
				}
				fmt.Print(" CPRUA:", cprua)
			}
			if sid == 31 {
				cprub := CPRU{}
				err = cprub.read(r)
				if err != nil {
					log.Fatal("cprub", err)
				}
				fmt.Print(" CPRUB:", cprub)
			}

		}

		fmt.Println()
	}

}
