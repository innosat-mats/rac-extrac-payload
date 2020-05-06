package extractors

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// DecodeAEZ parses AEZ packages
func DecodeAEZ(target chan<- common.DataRecord, source <-chan common.DataRecord) {
	defer close(target)
	for sourcePacket := range source {
		var exportable common.Exportable
		var err error
		buffer := bytes.NewBuffer(sourcePacket.Buffer)
		switch {
		case sourcePacket.TMHeader.IsHousekeeping():
			var sid aez.SID
			binary.Read(buffer, binary.BigEndian, &sid)
			sourcePacket.SID = sid
			exportable, err = instrumentHK(sid, buffer)
		case sourcePacket.TMHeader.IsTransparentData():
			var rid aez.RID
			binary.Read(buffer, binary.BigEndian, &rid)
			sourcePacket.RID = rid
			exportable, err = instrumentTransparentData(rid, buffer)
		default:
			err = errors.New("the TMHeader isn't recognized as either housekeeping or tranparent data")
		}
		sourcePacket.Error = err
		sourcePacket.Data = exportable
		sourcePacket.Buffer = buffer.Bytes()
		target <- sourcePacket
	}
}
