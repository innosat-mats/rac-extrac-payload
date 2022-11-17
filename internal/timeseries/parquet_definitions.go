package timeseries

import (
	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// CCDParquetRow holds the exportable parquet representation of a CCD record, including common attributes
type CCDParquetRow struct {
	common.DataRecordParquet
	aez.CCDImageParquet
}

// PMParquetRow holds the exportable parquet representation of a PM record, including common attributes
type PMParquetRow struct {
	common.DataRecordParquet
	aez.PMDataParquet
}

// HTRParquetRow holds the exportable parquet representation of a HTR record, including common attributes
type HTRParquetRow struct {
	common.DataRecordParquet
	aez.HTRParquet
}

// CPRUParquetRow holds the exportable parquet representation of a CPRU record, including common attributes
type CPRUParquetRow struct {
	common.DataRecordParquet
	aez.CPRUParquet
}

// PWRParquetRow holds the exportable parquet representation of a PWR record, including common attributes
type PWRParquetRow struct {
	common.DataRecordParquet
	aez.PWRParquet
}

// STATParquetRow holds the exportable parquet representation of a STAT record, including common attributes
type STATParquetRow struct {
	common.DataRecordParquet
	aez.STATParquet
}

// TCVParquetRow holds the exportable parquet representation of a TCV record, including common attributes
type TCVParquetRow struct {
	common.DataRecordParquet
	aez.TCVParquet
}

// GetParquetRow returns the exportable parquet representation of a record, including common attributes
func GetParquetRow(pkg *common.DataRecord) interface{} {
	switch pkg.Data.(type) {
	case *aez.CCDImage:
		ccd, ok := pkg.Data.(*aez.CCDImage)
		if ok {
			return CCDParquetRow{
				pkg.GetParquet(),
				ccd.GetParquet(),
			}
		}
		return CCDParquetRow{
			pkg.GetParquet(),
			aez.CCDImageParquet{},
		}
	case *aez.PMData:
		pm, ok := pkg.Data.(*aez.PMData)
		if ok {
			return PMParquetRow{
				pkg.GetParquet(),
				pm.GetParquet(),
			}
		}
		return PMParquetRow{
			pkg.GetParquet(),
			aez.PMDataParquet{},
		}
	case *aez.HTR:
		htr, ok := pkg.Data.(*aez.HTR)
		if ok {
			return HTRParquetRow{
				pkg.GetParquet(),
				htr.GetParquet(),
			}
		}
		return HTRParquetRow{
			pkg.GetParquet(),
			aez.HTRParquet{},
		}
	case *aez.PWR:
		pwr, ok := pkg.Data.(*aez.PWR)
		if ok {
			return PWRParquetRow{
				pkg.GetParquet(),
				pwr.GetParquet(),
			}
		}
		return PWRParquetRow{
			pkg.GetParquet(),
			aez.PWRParquet{},
		}
	case *aez.CPRU:
		cpru, ok := pkg.Data.(*aez.CPRU)
		if ok {
			return CPRUParquetRow{
				pkg.GetParquet(),
				cpru.GetParquet(),
			}
		}
		return CPRUParquetRow{
			pkg.GetParquet(),
			aez.CPRUParquet{},
		}
	case *aez.STAT:
		stat, ok := pkg.Data.(*aez.STAT)
		if ok {
			return STATParquetRow{
				pkg.GetParquet(),
				stat.GetParquet(),
			}
		}
		return STATParquetRow{
			pkg.GetParquet(),
			aez.STATParquet{},
		}
	case *aez.TCAcceptSuccessData:
		tcv, ok := pkg.Data.(*aez.TCAcceptSuccessData)
		if ok {
			return TCVParquetRow{
				pkg.GetParquet(),
				tcv.GetParquet(),
			}
		}
		return TCVParquetRow{
			pkg.GetParquet(),
			aez.TCVParquet{},
		}
	case *aez.TCAcceptFailureData:
		tcv, ok := pkg.Data.(*aez.TCAcceptFailureData)
		if ok {
			return TCVParquetRow{
				pkg.GetParquet(),
				tcv.GetParquet(),
			}
		}
		return TCVParquetRow{
			pkg.GetParquet(),
			aez.TCVParquet{},
		}
	case *aez.TCExecSuccessData:
		tcv, ok := pkg.Data.(*aez.TCExecSuccessData)
		if ok {
			return TCVParquetRow{
				pkg.GetParquet(),
				tcv.GetParquet(),
			}
		}
		return TCVParquetRow{
			pkg.GetParquet(),
			aez.TCVParquet{},
		}
	case *aez.TCExecFailureData:
		tcv, ok := pkg.Data.(*aez.TCExecFailureData)
		if ok {
			return TCVParquetRow{
				pkg.GetParquet(),
				tcv.GetParquet(),
			}
		}
		return TCVParquetRow{
			pkg.GetParquet(),
			aez.TCVParquet{},
		}
	}
	return pkg.GetParquet()
}
