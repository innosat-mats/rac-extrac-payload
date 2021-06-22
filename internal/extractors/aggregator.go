package extractors

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

// Aggregator sorts and accumulates standalone and multi-packets
func Aggregator(target chan<- common.DataRecord, source <-chan common.DataRecord) {
	defer close(target)
	multiPackBuffer := bytes.NewBuffer([]byte{})
	var multiPackStarted bool
	var multiPackStart common.DataRecord
	const sidRidLength = 2

	for sourcePacket := range source {
		if sourcePacket.Error != nil {
			target <- sourcePacket
			continue
		}
		switch sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags() {
		case innosat.SPStandalone:
			// Produce error for unfinished multipack lingering
			if multiPackStarted {
				target <- makeUnfinishedMultiPackError(multiPackBuffer, sourcePacket)
				multiPackBuffer = bytes.NewBuffer([]byte{})
				multiPackStarted = false
			}

			// Report standalone pack
			target <- sourcePacket
		case innosat.SPStart:
			// Produce error for unfinished multipack lingering
			if multiPackStarted {
				target <- makeUnfinishedMultiPackError(multiPackBuffer, sourcePacket)
			}

			// Start new multipack
			buffer := bytes.NewBuffer(sourcePacket.Buffer)
			multiPackStarted = true
			multiPackBuffer = bytes.NewBuffer([]byte{})
			multiPackStart = sourcePacket
			_, err := multiPackBuffer.ReadFrom(buffer)
			if err != nil && err != io.EOF {
				sourcePacket.Error = err
				target <- sourcePacket
			}
		case innosat.SPCont:
			// Report error missing start packet
			if !multiPackStarted {
				sourcePacket.Error = errors.New(
					"got continuation packet without a start packet",
				)
				multiPackStart = sourcePacket
				multiPackStarted = true
			}

			// Concat SPCont packet
			buffer := bytes.NewBuffer(sourcePacket.Buffer[sidRidLength:len(sourcePacket.Buffer)])
			_, err := multiPackBuffer.ReadFrom(buffer)
			if err != nil && err != io.EOF {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = err
				target <- sourcePacketCopy
			}
		case innosat.SPStop:
			// Report error missing start pack
			if !multiPackStarted {
				// Report error
				multiPackStart = sourcePacket
				multiPackStart.Error = errors.New(
					"got stop packet without a start packet",
				)
			}

			// Concat SPStop and report parsed packet
			buffer := bytes.NewBuffer(sourcePacket.Buffer[sidRidLength:len(sourcePacket.Buffer)])
			_, err := multiPackBuffer.ReadFrom(buffer)
			if err != nil && err != io.EOF {
				sourcePacket.Error = err
				target <- sourcePacket
			}
			multiPackStart.Buffer = multiPackBuffer.Bytes()
			target <- multiPackStart
			multiPackBuffer = bytes.NewBuffer([]byte{})
			multiPackStart = common.DataRecord{}
			multiPackStarted = false

		default:
			// Report unknown grouping flag error
			sourcePacket.Error = fmt.Errorf(
				"unhandled grouping flag %v",
				sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags(),
			)
			sourcePacket.Buffer = multiPackBuffer.Bytes()
			target <- sourcePacket
		}
	}

	// Report attemmpt at parsing dangling multipack
	if multiPackStarted {
		err := fmt.Errorf(
			"dangling final multipacket with %v bytes",
			multiPackBuffer.Len(),
		)
		multiPackStart.Buffer = multiPackBuffer.Bytes()
		multiPackStart.Error = err
		target <- multiPackStart
	}
}

func makeUnfinishedMultiPackError(multiPackBuffer *bytes.Buffer, sourcePacket common.DataRecord) common.DataRecord {
	errorPacket := sourcePacket
	errorPacket.Error = errors.New(
		"orphaned multi-package data without termination detected",
	)
	errorPacket.Buffer = multiPackBuffer.Bytes()
	return errorPacket
}
