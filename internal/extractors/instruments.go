package extractors

import (
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

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
		ccdImg := aez.CCDImage{PackData: ccdIPD, BadColumns: badColumns}
		dataPackage = ccdImg
	default:
		err = fmt.Errorf("unhandled RID %v", rid)
	}
	return dataPackage, err
}
