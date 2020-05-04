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

//PackageType interface for the data read from AEZ
type PackageType interface {
}

// DecodeAEZ processes data packages
func DecodeAEZ(target chan<- common.DataRecord, source <-chan common.DataRecord) {
	defer close(target)
	var multiPackBuffer *bytes.Buffer

	for sourcePacket := range source {
		if sourcePacket.Error != nil {
			target <- sourcePacket
			continue
		}
		reader := bytes.NewReader(sourcePacket.Buffer)
		switch sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags() {
		case innosat.SPStandalone:
			// Panic if has started multipack
			if multiPackBuffer != nil {
				target <- prepareUnfinishedMultiPackPanic(multiPackBuffer, sourcePacket)
				multiPackBuffer = nil
			}
			target <- processSourcePacket(reader, sourcePacket)
		case innosat.SPStart:
			// Panic if hasn't concluded previous
			if multiPackBuffer != nil {
				target <- prepareUnfinishedMultiPackPanic(multiPackBuffer, sourcePacket)
			}
			multiPackBuffer = bytes.NewBuffer([]byte{})
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacket.Error = err
				target <- sourcePacket
			}
		case innosat.SPCont:
			// Panic if hasn't started package
			if multiPackBuffer == nil {
				multiPackBuffer = bytes.NewBuffer([]byte{})
				sourcePacket.Error = errors.New("got continuation packet without a start packet")
				target <- sourcePacket
			}
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = err
				target <- sourcePacketCopy
			}
		case innosat.SPStop:
			if multiPackBuffer == nil {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = errors.New("got stop packed without a start packet")
				target <- sourcePacketCopy
				multiPackBuffer = bytes.NewBuffer([]byte{})
			}
			_, err := io.Copy(multiPackBuffer, reader)
			if err != nil {
				sourcePacketCopy := sourcePacket
				sourcePacketCopy.Error = err
				target <- sourcePacketCopy
			}
			reader = bytes.NewReader(multiPackBuffer.Bytes())
			target <- processSourcePacket(reader, sourcePacket)

		default:
			sourcePacket.Error = fmt.Errorf("unhandled grouping flag %v", sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags())
			target <- sourcePacket
		}
	}
	//Panic if unclosed multipack ?
}

func processSourcePacket(reader *bytes.Reader, sourcePacket common.DataRecord) common.DataRecord {
	switch {
	case sourcePacket.TMHeader.IsHousekeeping():
		var sid aez.SID
		binary.Read(reader, binary.BigEndian, &sid)
		targetPacket, err := instrumentHK(sid, reader)
		var buf []byte
		_, err2 := reader.Read(buf)
		if err2 != nil {
			if err == nil {
				err = err2
			} else {
				err = fmt.Errorf("%v | %v", err, err2)
			}
		}
		sourcePacket.Data = targetPacket
		sourcePacket.SID = sid
		sourcePacket.Error = err
		sourcePacket.Buffer = buf
	case sourcePacket.TMHeader.IsTransparentData():
		var rid aez.RID
		binary.Read(reader, binary.BigEndian, &rid)
		targetPacket, err := instrumentTransparentData(rid, reader)
		var buf []byte
		_, err2 := reader.Read(buf)
		if err2 != nil {
			if err == nil {
				err = err2
			} else {
				err = fmt.Errorf("%v | %v", err, err2)
			}
		}
		sourcePacket.Data = targetPacket
		sourcePacket.RID = rid
		sourcePacket.Error = err
		sourcePacket.Buffer = buf
	default:
		sourcePacket.Error = errors.New("the TMHeader isn't recognized as either housekeeping or tranparent data")
	}
	return sourcePacket
}

func prepareUnfinishedMultiPackPanic(multiPackBuffer *bytes.Buffer, sourcePacket common.DataRecord) common.DataRecord {
	panicReport := sourcePacket
	panicReport.Error = errors.New("orphaned multi-package data without termination detected")
	panicReport.Buffer = multiPackBuffer.Bytes()
	return panicReport
}

func instrumentHK(sid aez.SID, buf io.Reader) (common.Exportable, error) {
	var dataPackage common.Exportable
	var err error
	switch sid {
	case aez.SIDSTAT:
		stat := aez.STAT{}
		err = stat.Read(buf)
		dataPackage = stat
	case aez.SIDHTR:
		htr := aez.HTR{}
		err = htr.Read(buf)
		dataPackage = htr
	case aez.SIDPWR:
		pwr := aez.PWR{}
		err = pwr.Read(buf)
		dataPackage = pwr
	case aez.SIDCPRUA:
		cpru := aez.CPRU{}
		err = cpru.Read(buf)
		dataPackage = cpru
	case aez.SIDCPRUB:
		cpru := aez.CPRU{}
		err = cpru.Read(buf)
		dataPackage = cpru
	default:
		err = fmt.Errorf("unhandled SID %v", sid)
	}
	return dataPackage, err
}

func instrumentTransparentData(rid aez.RID, buf io.Reader) (common.Exportable, error) {
	var dataPackage common.Exportable
	var err error
	switch rid {
	case aez.CCD1, aez.CCD2, aez.CCD3, aez.CCD4, aez.CCD5, aez.CCD6, aez.CCD7:
		ccdIPD := aez.CCDImagePackData{}
		var badColumns []uint16
		badColumns, err = ccdIPD.Read(buf)
		ccdI := aez.CCDImage{PackData: ccdIPD, BadColumns: badColumns}
		dataPackage = ccdI
	default:
		err = fmt.Errorf("unhandled RID %v", rid)
	}
	return dataPackage, err
}
