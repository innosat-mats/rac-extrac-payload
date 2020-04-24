package common

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// Packet is complete Ramses packet
type Packet struct {
	Header       ramses.Ramses
	SecureHeader ramses.Secure
	Payload      []byte
}

// StreamBatch tells origin of batch
type StreamBatch struct {
	Buf    io.Reader
	Origin OriginDescription
}

//DecodeRamses reads Ramses packages from buffer
func DecodeRamses(recordChannel chan<- DataRecord, streamBatch ...StreamBatch) {
	defer close(recordChannel)
	var err error
	for _, stream := range streamBatch {
		for {
			header := ramses.Ramses{}
			err = header.Read(stream.Buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				recordChannel <- DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}
				break
			}

			if !header.Valid() {
				err := fmt.Errorf("Not a valid RAC-file")
				recordChannel <- DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}
				break
			}

			secureHeader := ramses.Secure{}
			if header.SecureTrans() {
				err = secureHeader.Read(stream.Buf)
				if err != nil {
					recordChannel <- DataRecord{Origin: stream.Origin, RamsesHeader: header, Error: err, Buffer: []byte{}}
					break
				}
				header.Length -= uint16(binary.Size(secureHeader))
			}

			payload := make([]byte, header.Length)
			_, err = stream.Buf.Read(payload)
			if err != nil {
				recordChannel <- DataRecord{Origin: stream.Origin, RamsesHeader: header, Error: err, Buffer: []byte{}}
				break
			}
			recordChannel <- DataRecord{Origin: stream.Origin, RamsesHeader: header, RamsesSecure: secureHeader, Error: nil, Buffer: payload}
		}
	}
}
