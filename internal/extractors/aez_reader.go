package extractors

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

// DecodeAEZ processes data packages
func DecodeAEZ(target chan<- common.DataRecord, source <-chan common.DataRecord) {
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
			target <- processSourcePacket(reader, sourcePacket)
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
			reader = bytes.NewReader(multiPackBuffer.Bytes())
			target <- processSourcePacket(reader, multiPackStart)

		default:
			// Report unknown grouping flag error
			sourcePacket.Error = fmt.Errorf("unhandled grouping flag %v", sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags())
			target <- sourcePacket
		}
	}

	// Report attemmpt at parsing dangling multipack
	err := fmt.Errorf("dangling final multipacket with %v bytes", multiPackBuffer.Len())
	reader := bytes.NewReader(multiPackBuffer.Bytes())
	sourcePacket := processSourcePacket(reader, multiPackStart)
	sourcePacket.Error = err
	target <- sourcePacket
}

func processSourcePacket(reader *bytes.Reader, sourcePacket common.DataRecord) common.DataRecord {
	switch {
	case sourcePacket.TMHeader.IsHousekeeping():
		var sid aez.SID
		binary.Read(reader, binary.BigEndian, &sid)
		sourcePacket.SID = sid
		exportable, err := instrumentHK(sid, reader)
		addExportable(&sourcePacket, reader, &exportable, err)
	case sourcePacket.TMHeader.IsTransparentData():
		var rid aez.RID
		binary.Read(reader, binary.BigEndian, &rid)
		sourcePacket.RID = rid
		exportable, err := instrumentTransparentData(rid, reader)
		addExportable(&sourcePacket, reader, &exportable, err)
	default:
		sourcePacket.Error = errors.New("the TMHeader isn't recognized as either housekeeping or tranparent data")
	}
	return sourcePacket
}

func addExportable(sourcePacket *common.DataRecord, reader *bytes.Reader, exportable *common.Exportable, exportableErr error) {
	var buf []byte
	_, bufErr := reader.Read(buf)
	if bufErr != nil && bufErr != io.EOF {
		if exportableErr == nil {
			sourcePacket.Error = bufErr
		} else {
			sourcePacket.Error = fmt.Errorf("%v | %v", exportableErr, bufErr)
		}
	} else {
		sourcePacket.Error = exportableErr
	}
	sourcePacket.Data = *exportable
	sourcePacket.Buffer = buf

}

func makeUnfinishedMultiPackError(multiPackBuffer *bytes.Buffer, sourcePacket common.DataRecord) common.DataRecord {
	panicReport := sourcePacket
	panicReport.Error = errors.New("orphaned multi-package data without termination detected")
	panicReport.Buffer = multiPackBuffer.Bytes()
	return panicReport
}
