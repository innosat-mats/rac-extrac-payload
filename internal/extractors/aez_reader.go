package extractors

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// DecodeAEZ parses AEZ packages
func DecodeAEZ(target chan<- common.DataRecord, source <-chan common.DataRecord) {
	defer close(target)
	var exportable common.Exportable
	var err error
	var buffer *bytes.Buffer
	for sourcePacket := range source {
		if sourcePacket.Error != nil {
			target <- sourcePacket
			continue
		}
		buffer = bytes.NewBuffer(sourcePacket.Buffer)
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
		case sourcePacket.TMHeader.IsTCVerification():
			exportable, err = instrumentVerification(sourcePacket.TMHeader.ServiceSubType, buffer)
		default:
			err = errors.New("the TMHeader isn't recognized as either housekeeping, transparent or verification data")
			exportable = nil
		}
		if err != io.EOF {
			sourcePacket.Error = err
		}
		sourcePacket.Data = exportable
		sourcePacket.Buffer = buffer.Bytes()
		target <- sourcePacket
	}
}
