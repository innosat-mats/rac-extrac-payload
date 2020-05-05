package extractors

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// DecodeAEZ parses AEZ packages
func DecodeAEZ(target chan<- common.DataRecord, source <-chan common.DataRecord) {
	defer close(target)
	for sourcePacket := range source {
		reader := bytes.NewReader(sourcePacket.Buffer)
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
		target <- sourcePacket
	}
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
