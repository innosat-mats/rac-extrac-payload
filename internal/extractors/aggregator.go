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

	for sourcePacket := range source {
		if sourcePacket.Error != nil {
			target <- sourcePacket
			continue
		}
		reader := bytes.NewReader(sourcePacket.Buffer)
		switch sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags() {
		case innosat.SPStandalone:
			// Produce error for unfinished multipack lingering
			if multiPackStarted {
				target <- makeUnfinishedMultiPackError(multiPackBuffer, sourcePacket)
				multiPackBuffer.Reset()
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
			multiPackStarted = true
			multiPackBuffer.Reset()
			multiPackStart = sourcePacket
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil && err != io.EOF {
				sourcePacket.Error = err
				target <- sourcePacket
			}
		case innosat.SPCont:
			// Report error missing start packet
			if !multiPackStarted {
				sourcePacket.Error = errors.New("got continuation packet without a start packet")
				target <- sourcePacket

				multiPackStart = sourcePacket
				multiPackStarted = true
			}

			// Concat SPCont packet
			_, err := multiPackBuffer.ReadFrom(reader)
			if err != nil && err != io.EOF {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = err
				target <- sourcePacketCopy
			}
		case innosat.SPStop:
			// Report error missing start pack
			if !multiPackStarted {
				// Report error
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = errors.New("got stop packet without a start packet")
				target <- sourcePacketCopy

				multiPackStart = sourcePacket
			}

			// Concat SPStop and report parsed packet
			_, err := multiPackBuffer.ReadFrom(reader)
			if err != nil && err != io.EOF {
				sourcePacket.Error = err
				target <- sourcePacket
			}
			multiPackStart.Buffer = multiPackBuffer.Bytes()
			target <- multiPackStart
			multiPackBuffer.Reset()
			multiPackStart = common.DataRecord{}
			multiPackStarted = false

		default:
			// Report unknown grouping flag error
			sourcePacket.Error = fmt.Errorf("unhandled grouping flag %v", sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags())
			target <- sourcePacket
		}
	}

	// Report attemmpt at parsing dangling multipack
	if multiPackStarted {
		err := fmt.Errorf("dangling final multipacket with %v bytes", multiPackBuffer.Len())
		multiPackStart.Buffer = multiPackBuffer.Bytes()
		multiPackStart.Error = err
		target <- multiPackStart
	}
}

func makeUnfinishedMultiPackError(multiPackBuffer *bytes.Buffer, sourcePacket common.DataRecord) common.DataRecord {
	errorPacket := sourcePacket
	errorPacket.Error = errors.New("orphaned multi-package data without termination detected")
	errorPacket.Buffer = multiPackBuffer.Bytes()
	return errorPacket
}
