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
	Header   ramses.Ramses
	TMHeader ramses.TMHeader
	Payload  []byte
}

// StreamBatch tells origin of batch
type StreamBatch struct {
	Buf    io.Reader
	Origin common.OriginDescription
}

//DecodeRamses reads Ramses packages from buffer
func DecodeRamses(recordChannel chan<- common.DataRecord, streamBatch ...StreamBatch) {
	defer close(recordChannel)
	var records []common.DataRecord

	for _, stream := range streamBatch {
		record, done := getRecord(stream)
		if done {
			continue
		}
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
				record, done := getRecord(stream)
				if done {
					break
				}
				recordChannel <- record
				if record.Error != nil {
					break
				}
			}
		}
	}
}

// getRecord returns next record and a flag if stream was actually done prior to this record
func getRecord(stream StreamBatch) (common.DataRecord, bool) {
	header, err := ramses.NewRamses(stream.Buf)
	if err != nil {
		// EOF for reading Ramses just means we have no more records so not really an error
		if err == io.EOF {
			return common.DataRecord{}, true
		}
		return common.DataRecord{
			Origin: stream.Origin,
			Error:  fmt.Errorf("could not parse ramses header: %v (%v)", err, stream.Origin.Name),
			Buffer: []byte{},
		}, false
	}

	if !header.Valid() {
		err := fmt.Errorf("Not a valid RAC-record %v (%s)", header, stream.Origin.Name)
		return common.DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}, false
	}

	tmHeader, err := ramses.NewTMHeader(stream.Buf)
	if err != nil && err != io.EOF {
		return common.DataRecord{
			Origin:       stream.Origin,
			RamsesHeader: header,
			Error:        fmt.Errorf("could not parse OHBSE CCDS TM Packet header: %v (%v)", err, stream.Origin.Name),
			Buffer:       []byte{},
		}, false
	}
	header.Length -= uint16(binary.Size(tmHeader))

	payload := make([]byte, header.Length)
	n, _ := stream.Buf.Read(payload)
	if n != int(header.Length) {
		return common.DataRecord{
			Origin:         stream.Origin,
			RamsesHeader:   header,
			RamsesTMHeader: tmHeader,
			Error: fmt.Errorf(
				"payload truncaded, only found %v bytes but needed %v (%v)",
				n,
				header.Length,
				stream.Origin.Name,
			),
			Buffer: []byte{},
		}, false
	}

	return common.DataRecord{
		Origin:         stream.Origin,
		RamsesHeader:   header,
		RamsesTMHeader: tmHeader,
		Error:          nil,
		Buffer:         payload,
	}, false
}
