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
		dataPackage, err = aez.NewSTAT(buf)
	case aez.SIDHTR:
		dataPackage, err = aez.NewHTR(buf)
	case aez.SIDPWR:
		dataPackage, err = aez.NewPWR(buf)
	case aez.SIDCPRUA, aez.SIDCPRUB:
		dataPackage, err = aez.NewCPRU(buf)
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
		dataPackage, err = aez.NewPMData(buf)
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
		dataPackage, err = aez.NewTCAcceptSuccessData(buf)
	case innosat.TCAcceptFailure:
		dataPackage, err = aez.NewTCAcceptFailureData(buf)
	case innosat.TCExecSuccess:
		dataPackage, err = aez.NewTCExecSuccessData(buf)
	case innosat.TCExecFailure:
		dataPackage, err = aez.NewTCExecFailureData(buf)
	default:
		err = fmt.Errorf("unhandled TC Verification subtype %v", subtype)
	}
	return dataPackage, err
}
