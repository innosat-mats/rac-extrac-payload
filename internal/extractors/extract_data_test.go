package extractors

import (
	"bytes"
	"fmt"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

// Example ...
func Example() {

	data := []byte{
		//1 A continuation packet (started in previous file), should error in aggregator
		// Ramses header
		0x90, 0xeb, 0x58, 0x00, 0x79, 0xd8, 0x00, 0x00,
		0xdc, 0xad, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00,
		// Ramses TM header
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00,
		0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		// Ramses TM Payload
		0x08, 0x64, 0x88, 0x97, 0x00, 0x41, 0x10, 0x80,
		0x19, 0x00, 0x00, 0x12, 0x19, 0xda, 0x7e, 0x00,
		0x17, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x2f, 0x01,
		0x2f, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x2e, 0x01,
		0x32, 0x01, 0x30, 0x01, 0x2f, 0x01, 0x2f, 0x01,
		0x2f, 0x01, 0x2d, 0x01, 0x30, 0x01, 0x2e, 0x01,
		0x30, 0x01, 0x2b, 0x01, 0x2f, 0x01, 0x31, 0x01,
		0x32, 0x01, 0x32, 0x01, 0x33, 0x01, 0x2e, 0x01,
		0x31, 0x01, 0x2d, 0x01, 0x2e, 0x01, 0x20, 0xf7,

		//2 A STAT package should parse all the way through
		// Ramses header
		0x90, 0xeb, 0x48, 0x00, 0x79, 0xd8, 0x00, 0x00,
		0xfb, 0xad, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00,
		// Ramses TM header
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00,
		0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		// Ramses TM Payload
		0x08, 0x64, 0xc8, 0x98, 0x00, 0x31, 0x10, 0x03,
		0x19, 0x00, 0x00, 0x12, 0x19, 0xe3, 0x39, 0x00,
		0x01, 0x7f, 0x04, 0x02, 0x82, 0x04, 0x02, 0x02,
		0x06, 0x01, 0x19, 0x12, 0x00, 0x00, 0x0c, 0xe3,
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0xfb,

		//3 Too short buffer (innosat tm header length too short for a STAT) and also bad CRC
		// Package length for ramses is correct (3rd & 4th byte here 0x46 = 70, and the 16 bytes header)
		// Ramses header
		0x90, 0xeb, 0x46, 0x00, 0x79, 0xd8, 0x00, 0x00,
		0xd9, 0xb0, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00,
		// Ramses TM header
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00,
		0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		// Ramses TM Payload
		0x08, 0x64, 0xc8, 0x99, 0x00, 0x31, 0x10, 0x03,
		0x19, 0x00, 0x00, 0x12, 0x1a, 0xa0, 0xec, 0x00,
		0x01, 0x7f, 0x04, 0x02, 0x82, 0x04, 0x02, 0x02,
		0x06, 0x01, 0x1a, 0x12, 0x00, 0x00, 0xd0, 0xa0,
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x12, 0x48,

		//4 A STAT package should parse through
		// Ramses header
		0x90, 0xeb, 0x48, 0x00, 0x79, 0xd8, 0x00, 0x00,
		0xc1, 0xb4, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00,
		// Ramses TM header
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00,
		0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		// Ramses TM Payload
		0x08, 0x64, 0xc8, 0x9a, 0x00, 0x31, 0x10, 0x03,
		0x19, 0x00, 0x00, 0x12, 0x1b, 0xa1, 0x0b, 0x00,
		0x01, 0x7f, 0x04, 0x02, 0x82, 0x04, 0x02, 0x02,
		0x06, 0x01, 0x1b, 0x12, 0x00, 0x00, 0xef, 0xa0,
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa5, 0xd5,

		//5 Checksum bad (CRC)
		// Ramses header
		0x90, 0xeb, 0x48, 0x00, 0x79, 0xd8, 0x00, 0x00,
		0xa9, 0xb8, 0xc8, 0x04, 0xda, 0x1c, 0x00, 0x00,
		// Ramses TM header
		0x00, 0x00, 0x2e, 0x02, 0x00, 0x00, 0x64, 0x00,
		0x00, 0x00, 0x00, 0xcc, 0xcc, 0xcc, 0xcc, 0x00,
		// Ramses TM Payload
		0x08, 0x64, 0xc8, 0x9b, 0x00, 0x31, 0x10, 0x03,
		0x19, 0x00, 0x00, 0x12, 0x1c, 0xa0, 0xe8, 0x00,
		0x01, 0x7f, 0x04, 0x02, 0x82, 0x04, 0x02, 0x02,
		0x06, 0x01, 0x1c, 0x12, 0x00, 0x00, 0xcd, 0xa0,
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7e, 0x00,
	}
	reader1 := bytes.NewReader(data)
	reader2 := bytes.NewReader(data)
	ExtractData(
		simpleOutput,
		Slask{},
		StreamBatch{reader1, &common.OriginDescription{Name: "Set1", ProcessingDate: innosat.Epoch}},
		StreamBatch{reader2, &common.OriginDescription{Name: "Set2", ProcessingDate: innosat.Epoch}},
	)

	// Output:
	// got stop packet without a start packet
	//  STAT <nil>
	//   checksum bad 62691
	//  STAT <nil>
	//   checksum bad 32488
	//   got stop packet without a start packet
	//  STAT <nil>
	//   checksum bad 62691
	//  STAT <nil>
	//   checksum bad 32488
}

func simpleOutput(pkg common.DataRecord) {
	fmt.Println(pkg.RID.String(), pkg.SID.String(), pkg.Error)
}
