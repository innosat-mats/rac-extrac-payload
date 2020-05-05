package extractors

import (
	"bytes"
	"encoding/binary"
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
	for sourcePacket := range source {
		if sourcePacket.Error != nil {
			target <- sourcePacket
			continue
		}
		reader := bytes.NewReader(sourcePacket.Buffer)
		if sourcePacket.SourceHeader.PacketSequenceControl.GroupingFlags() == innosat.SPStandalone {
			switch {
			case sourcePacket.TMHeader.IsHousekeeping():
				var sid aez.SID
				binary.Read(reader, binary.BigEndian, &sid)
				targetPacket, err := instrumentHK(sid, reader)
				sourcePacket.Data = targetPacket
				sourcePacket.SID = sid
				sourcePacket.Error = err
				if reader.Len() == 0 {
					sourcePacket.Buffer = []byte{}
				}
				target <- sourcePacket
			}
		}
	}
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
	}
	return dataPackage, err
}
