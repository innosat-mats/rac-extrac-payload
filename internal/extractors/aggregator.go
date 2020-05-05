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
	var multiPackBuffer *bytes.Buffer
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
			if multiPackBuffer != nil {
				target <- makeUnfinishedMultiPackError(multiPackBuffer, sourcePacket)
				multiPackBuffer = nil
			}

			// Report standalone pack
			target <- sourcePacket
		case innosat.SPStart:
			// Produce error for unfinished multipack lingering
			if multiPackBuffer != nil {
				target <- makeUnfinishedMultiPackError(multiPackBuffer, sourcePacket)
			}

			// Start new multipack
			multiPackBuffer = bytes.NewBuffer([]byte{})
			multiPackStart = sourcePacket
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacket.Error = err
				target <- sourcePacket
			}
		case innosat.SPCont:
			// Report error missing start packet
			if multiPackBuffer == nil {
				sourcePacket.Error = errors.New("got continuation packet without a start packet")
				target <- sourcePacket
				// Create a new one to capture remaining SPCont and attempt parse
				multiPackBuffer = bytes.NewBuffer([]byte{})
				multiPackStart = sourcePacket
			}

			// Concat SPCont packet
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = err
				target <- sourcePacketCopy
			}
		case innosat.SPStop:
			// Report error missing start pack
			if multiPackBuffer == nil {
				// Create a new one to capture the output and attempt to parse
				multiPackBuffer = bytes.NewBuffer([]byte{})
				multiPackStart = sourcePacket

				// Report error
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = errors.New("got stop packed without a start packet")
				target <- sourcePacketCopy
			}

			// Concat SPStop and report parsed packet
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacket.Error = err
				target <- sourcePacket
			}
			multiPackStart.Buffer = multiPackBuffer.Bytes()
			target <- multiPackStart
			multiPackBuffer = nil
			multiPackStart = common.DataRecord{}

		default:
			// Report unknown grouping flag error
			sourcePacket.Error = fmt.Errorf("unhandled grouping flag %v", sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags())
			target <- sourcePacket
		}
	}

	// Report attemmpt at parsing dangling multipack
	err := fmt.Errorf("dangling final multipacket with %v bytes", multiPackBuffer.Len())
	multiPackStart.Buffer = multiPackBuffer.Bytes()
	multiPackStart.Error = err
	target <- multiPackStart
}

func makeUnfinishedMultiPackError(multiPackBuffer *bytes.Buffer, sourcePacket common.DataRecord) common.DataRecord {
	errorPacket := sourcePacket
	errorPacket.Error = errors.New("orphaned multi-package data without termination detected")
	errorPacket.Buffer = multiPackBuffer.Bytes()
	return errorPacket
}
