package parquetrow

import (
	"time"
)

// ParquetRow holds the exportable parquet representation of a CCD record, including common attributes
type ParquetRow struct {
	OriginFile          string    `parquet:"OriginFile"`
	ProcessingTime      time.Time `parquet:"ProcessingTime"`
	RamsesTime          time.Time `parquet:"RamsesTime"`
	QualityIndicator    uint8     `parquet:"QualityIndicator"`
	LossFlag            uint8     `parquet:"LossFlag"`
	VCFrameCounter      uint8     `parquet:"VCFrameCounter"`
	SPSequenceCount     uint16    `parquet:"SPSequenceCount"`
	TMHeaderTime        time.Time `parquet:"TMHeaderTime"`
	TMHeaderNanoseconds int64     `parquet:"TMHeaderNanoseconds"`
	SID                 string    `parquet:"SID"`
	RID                 string    `parquet:"RID"`

	CCDSEL             uint8     `parquet:"CCDSEL"`
	EXPNanoseconds     int64     `parquet:"EXPNanoseconds"`
	EXPDate            time.Time `parquet:"EXPDate"`
	WDWMode            string    `parquet:"WDWMode"`
	WDWInputDataWindow string    `parquet:"WDWInputDataWindow"`
	WDWOV              uint16    `parquet:"WDWOV"`
	JPEGQ              uint8     `parquet:"JPEGQ"`
	FRAME              uint16    `parquet:"FRAME"`
	NROW               uint16    `parquet:"NROW"`
	NRBIN              uint16    `parquet:"NRBIN"`
	NRSKIP             uint16    `parquet:"NRSKIP"`
	NCOL               uint16    `parquet:"NCOL"`
	NCBINFPGAColumns   int       `parquet:"NCBINFPGAColumns"`
	NCBINCCDColumns    int       `parquet:"NCBINCCDColumns"`
	NCSKIP             uint16    `parquet:"NCSKIP"`
	NFLUSH             uint16    `parquet:"NFLUSH"`
	TEXPMS             uint32    `parquet:"TEXPMS"`
	GAINMode           string    `parquet:"GAINMode"`
	GAINTiming         string    `parquet:"GAINTiming"`
	GAINTruncation     uint8     `parquet:"GAINTruncation"`
	TEMP               uint16    `parquet:"TEMP"`
	FBINOV             uint16    `parquet:"FBINOV"`
	LBLNK              uint16    `parquet:"LBLNK"`
	TBLNK              uint16    `parquet:"TBLNK"`
	ZERO               uint16    `parquet:"ZERO"`
	TIMING1            uint16    `parquet:"TIMING1"`
	TIMING2            uint16    `parquet:"TIMING2"`
	VERSION            uint16    `parquet:"VERSION"`
	TIMING3            uint16    `parquet:"TIMING3"`
	NBC                uint16    `parquet:"NBC"`
	BC                 []uint16  `parquet:"BadColumns"`
	ImageName          string    `parquet:"ImageName"`
	ImageData          []byte    `parquet:"ImageData"`

	PMTime        time.Time `parquet:"PMTime"`
	PMNanoseconds int64     `parquet:"PMNanoseconds"`
	PM1A          uint32    `parquet:"PM1A"`
	PM1ACNTR      uint32    `parquet:"PM1ACNTR"`
	PM1B          uint32    `parquet:"PM1B"`
	PM1BCNTR      uint32    `parquet:"PM1BCNTR"`
	PM1S          uint32    `parquet:"PM1S"`
	PM1SCNTR      uint32    `parquet:"PM1SCNTR"`
	PM2A          uint32    `parquet:"PM2A"`
	PM2ACNTR      uint32    `parquet:"PM2ACNTR"`
	PM2B          uint32    `parquet:"PM2B"`
	PM2BCNTR      uint32    `parquet:"PM2BCNTR"`
	PM2S          uint32    `parquet:"PM2S"`
	PM2SCNTR      uint32    `parquet:"PM2SCNTR"`

	HTR1A  float64 `parquet:"HTR1A"`
	HTR1B  float64 `parquet:"HTR1B"`
	HTR1OD float64 `parquet:"HTR1OD"`
	HTR2A  float64 `parquet:"HTR2A"`
	HTR2B  float64 `parquet:"HTR2B"`
	HTR2OD float64 `parquet:"HTR2OD"`
	HTR7A  float64 `parquet:"HTR7A"`
	HTR7B  float64 `parquet:"HTR7B"`
	HTR7OD float64 `parquet:"HTR7OD"`
	HTR8A  float64 `parquet:"HTR8A"`
	HTR8B  float64 `parquet:"HTR8B"`
	HTR8OD float64 `parquet:"HTR8OD"`

	PWRT    float64 `parquet:"PWRT"`    // Temp. sense ⁰C
	PWRP32V float64 `parquet:"PWRP32V"` // +32V voltage sense voltage
	PWRP32C float64 `parquet:"PWRP32C"` // +32V current sense current
	PWRP16V float64 `parquet:"PWRP16V"` // +16V voltage sense voltage
	PWRP16C float64 `parquet:"PWRP16C"` // +16V current sense current
	PWRM16V float64 `parquet:"PWRM16V"` // -16V voltage sense voltage
	PWRM16C float64 `parquet:"PWRM16C"` // -16V current sense current
	PWRP3V3 float64 `parquet:"PWRP3V3"` // +3V3 voltage sense voltage
	PWRP3C3 float64 `parquet:"PWRP3C3"` // +3V3 current sense current

	VGATE0       float64 `parquet:"VGATE0"`
	VSUBS0       float64 `parquet:"VSUBS0"`
	VRD0         float64 `parquet:"VRD0"`
	VOD0         float64 `parquet:"VOD0"`
	Overvoltage0 bool    `parquet:"Overvoltage0"`
	Power0       bool    `parquet:"Power0"`
	VGATE1       float64 `parquet:"VGATE1"`
	VSUBS1       float64 `parquet:"VSUBS1"`
	VRD1         float64 `parquet:"VRD1"`
	VOD1         float64 `parquet:"VOD1"`
	Overvoltage1 bool    `parquet:"Overvoltage1"`
	Power1       bool    `parquet:"Power1"`
	VGATE2       float64 `parquet:"VGATE2"`
	VSUBS2       float64 `parquet:"VSUBS2"`
	VRD2         float64 `parquet:"VRD2"`
	VOD2         float64 `parquet:"VOD2"`
	Overvoltage2 bool    `parquet:"Overvoltage2"`
	Power2       bool    `parquet:"Power2"`
	VGATE3       float64 `parquet:"VGATE3"`
	VSUBS3       float64 `parquet:"VSUBS3"`
	VRD3         float64 `parquet:"VRD3"`
	VOD3         float64 `parquet:"VOD3"`
	Overvoltage3 bool    `parquet:"Overvoltage3"`
	Power3       bool    `parquet:"Power3"`

	STATTime        time.Time `parquet:"STATTime"`
	STATNanoseconds int64     `parquet:"STATNanoseconds"`
	SPID            uint16    `parquet:"SPID"`
	SPREV           uint8     `parquet:"SPREV"`
	FPID            uint16    `parquet:"FPID"`
	FPREV           uint8     `parquet:"FPREV"`
	SVNA            uint8     `parquet:"SVNA"`
	SVNB            uint8     `parquet:"SVNB"`
	SVNC            uint8     `parquet:"SVNC"`
	MODE            uint8     `parquet:"MODE"`
	EDACE           uint32    `parquet:"EDACE"`
	EDACCE          uint32    `parquet:"EDACCE"`
	EDACN           uint32    `parquet:"EDACN"`
	SPWEOP          uint32    `parquet:"SPWEOP"`
	SPWEEP          uint32    `parquet:"SPWEEP"`
	ANOMALY         uint8     `parquet:"ANOMALY"`

	TCV       string `parquet:"TCV"`
	TCPID     uint16 `parquet:"TCPID"`
	PSC       uint16 `parquet:"PSC"`
	ErrorCode uint8  `parquet:"ErrorCode"`

	Warnings []string `parquet:"Warnings"`
	Errors   []string `parquet:"Errors"`
}
