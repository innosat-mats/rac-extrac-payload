package extractors

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// Packet is complete Ramses packet
type Packet struct {
	Header             ramses.Ramses
	OhbseCcsdsTMPacket ramses.OhbseCcsdsTMPacket
	Payload            []byte
}

// StreamBatch tells origin of batch
type StreamBatch struct {
	Buf    *common.RemainingReader
	Origin common.OriginDescription
}

//DecodeRamses reads Ramses packages from buffer
func DecodeRamses(recordChannel chan<- common.DataRecord, streamBatch ...StreamBatch) {
	defer close(recordChannel)
	var records []common.DataRecord

	for _, stream := range streamBatch {
		record := getRecord(stream)
		if record.Error != nil {
			recordChannel <- record
		} else {
			records = append(records, record)
		}
	}

	sort.SliceStable(
		records,
		func(i int, j int) bool {
			return records[i].RamsesHeader.Nanoseconds() < records[j].RamsesHeader.Nanoseconds()
		},
	)

	for _, firstRecord := range records {
		recordChannel <- firstRecord
		for _, stream := range streamBatch {
			if stream.Origin.Name != firstRecord.Origin.Name {
				continue
			}
			for {
				if stream.Buf.Len() == 0 {
					break
				}
				recordChannel <- getRecord(stream)
			}
		}
	}
}

func getRecord(stream StreamBatch) common.DataRecord {
	header := ramses.Ramses{}
	err := header.Read(stream.Buf)
	if err != nil {
		return common.DataRecord{
			Origin: stream.Origin,
			Error:  fmt.Errorf("could not parse ramses header: %v (%v)", err, stream.Origin.Name),
			Buffer: []byte{},
		}
	}

	if !header.Valid() {
		err := fmt.Errorf("Not a valid RAC-record (%v)", stream.Origin.Name)
		return common.DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}
	}

	ccsdsTMPacketHeader := ramses.OhbseCcsdsTMPacket{}
	err = ccsdsTMPacketHeader.Read(stream.Buf)
	if err != nil && err != io.EOF {
		return common.DataRecord{
			Origin:       stream.Origin,
			RamsesHeader: header,
			Error:        fmt.Errorf("could not parse OHBSE CCDS TM Packet header: %v (%v)", err, stream.Origin.Name),
			Buffer:       []byte{},
		}
	}
	header.Length -= uint16(binary.Size(ccsdsTMPacketHeader))

	payload := make([]byte, header.Length)
	n, _ := stream.Buf.Read(payload)
	if n != int(header.Length) {
		return common.DataRecord{
			Origin:       stream.Origin,
			RamsesHeader: header,
			Error: fmt.Errorf(
				"payload truncaded, only found %v bytes but needed %v (%v)",
				n,
				header.Length,
				stream.Origin.Name,
			),
			Buffer: []byte{},
		}
	}

	return common.DataRecord{
		Origin:                   stream.Origin,
		RamsesHeader:             header,
		OhbseCcsdsTMPacketHeader: ccsdsTMPacketHeader,
		Error:                    nil,
		Buffer:                   payload,
	}
}
