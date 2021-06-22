package common

import (
	"fmt"
	"strings"
)

// MakePackageInfo ...
func MakePackageInfo(sourcePacket *DataRecord) string {
	infos := make([]string, 0)
	infos = append(infos, sourcePacket.OriginName())
	if sourcePacket.SourceHeader != nil {
		infos = append(
			infos,
			fmt.Sprintf("Packet ID %v", sourcePacket.SourceHeader.PacketID),
		)
	}
	if sourcePacket.RamsesTMHeader != nil {
		infos = append(
			infos,
			fmt.Sprintf("VC Frame Counter %v", sourcePacket.RamsesTMHeader.VCFrameCounter),
		)
	}
	if sourcePacket.RamsesHeader != nil {
		infos = append(
			infos,
			fmt.Sprintf(
				"Date %v, Time %v",
				sourcePacket.RamsesHeader.Date,
				sourcePacket.RamsesHeader.Time,
			),
		)
	}
	return fmt.Sprintf(
		"[%s]",
		strings.Join(infos[:], " / "),
	)
}
