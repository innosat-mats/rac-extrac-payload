package extractors

import (
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func instrumentHK(sid aez.SID, buf io.Reader) (common.Exporter, error) {
	var dataPackage common.Exporter
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

func instrumentTransparentData(rid aez.RID, buf io.Reader) (common.Exporter, error) {
	var dataPackage common.Exporter
	var err error
	switch {
	case rid.IsCCD():
		dataPackage, err = aez.NewCCDImage(buf)
	case rid == aez.PM:
		pmData := aez.PMData{}
		err = pmData.Read(buf)
		dataPackage = pmData
	default:
		err = fmt.Errorf("unhandled RID %v", rid)
	}
	return dataPackage, err
}

func instrumentVerification(
	subtype innosat.SourcePackageServiceSubtype,
	buf io.Reader,
) (common.Exporter, error) {
	var dataPackage common.Exporter
	var err error
	switch subtype {
	case innosat.TCAcceptSuccess:
		tcv := aez.TCAcceptSuccessData{}
		err = tcv.Read(buf)
		dataPackage = tcv
	case innosat.TCAcceptFailure:
		tcv := aez.TCAcceptFailureData{}
		err = tcv.Read(buf)
		dataPackage = tcv
	case innosat.TCExecSuccess:
		tcv := aez.TCExecSuccessData{}
		err = tcv.Read(buf)
		dataPackage = tcv
	case innosat.TCExecFailure:
		tcv := aez.TCExecFailureData{}
		err = tcv.Read(buf)
		dataPackage = tcv
	default:
		err = fmt.Errorf("unhandled TC Verification subtype %v", subtype)
	}
	return dataPackage, err
}
