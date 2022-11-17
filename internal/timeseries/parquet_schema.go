package timeseries

// RacSchema is the parquet schema for saving RAC data, one row per packet
const RacSchema = `message schema {
	required string OriginFile
	requried int64  ProcessingDate (TIMESPAMP(NANOS, true))
	required int64  RamsesTime (TIMESPAMP(NANOS, true))
	required uint8  QualityIndicator
	required uint8  LossFlag
	required uint8  VCFrameCounter
	required uint16 SPSequenceCount
	required int64  TMHeaderTime (TIMESTAMP(NANOS, true))
	required int64  TMHeaderNanoseconds
	required string SID
	required string RID

	optional uint8  CCDSEL
	optional int64  EXPNanoseconds
	optional int64  EXPDate (TIMESTAMP(NANOS, true))
	optional string WDWMode
	optional string WDWInputDataWindow
	optional uint16 WDWOV
	optional uint8  JPEGQ
	optional uint16 FRAME
	optional uint16 NROW
	optional uint16 NRBIN
	optional uint16 NRSKIP
	optional uint16 NCOL
	optional int    NCBINFPGAColumns
	optional int    NCBINCCDColumns
	optional uint16 NCSKIP
	optional uint16 NFLUSH
	optional uint32 TEXPMS
	optional string GAINMode
	optional string GAINTiming
	optional uint8  GAINTruncation
	optional uint16 TEMP
	optional uint16 FBINOV
	optional uint16 LBLNK
	optional uint16 TBLNK
	optional uint16 ZERO
	optional uint16 TIMING1
	optional uint16 TIMING2
	optional uint16 VERSION
	optional uint16 TIMING3
	optional uint16 NBC
	optional group  BadColumns (LIST) {
		repeated group list {
			required uint16 column
		}
	}
	optional string ImageFileName
	optional binary ImageFile

	optional int64  PMTime (TIMESTAMP(NANOS, true))
	optional int64  PMNanoseconds
	optional uint32 PM1A
	optional uint32 PM1ACNTR
	optional uint32 PM1B
	optional uint32 PM1BCNTR
	optional uint32 PM1S
	optional uint32 PM1SCNTR
	optional uint32 PM2A
	optional uint32 PM2ACNTR
	optional uint32 PM2B
	optional uint32 PM2BCNTR
	optional uint32 PM2S
	optional uint32 PM2SCNTR

	optional float64 HTR1A
	optional float64 HTR1B
	optional float64 HTR1OD
	optional float64 HTR2A
	optional float64 HTR2B
	optional float64 HTR2OD
	optional float64 HTR7A
	optional float64 HTR7B
	optional float64 HTR7OD
	optional float64 HTR8A
	optional float64 HTR8B
	optional float64 HTR8OD

	optional float64 PWRT
	optional float64 PWRP32V
	optional float64 PWRP32C
	optional float64 PWRP16V
	optional float64 PWRP16C
	optional float64 PWRM16V
	optional float64 PWRM16C
	optional float64 PWRP3V3
	optional float64 PWRP3C3

	optional float64 VGATE0
	optional float64 VSUBS0
	optional float64 VRD0
	optional float64 VOD0
	optional boolean Overvoltage0
	optional boolean Power0
	optional float64 VGATE1
	optional float64 VSUBS1
	optional float64 VRD1
	optional float64 VOD1
	optional boolean Overvoltage1
	optional boolean Power1
	optional float64 VGATE2
	optional float64 VSUBS2
	optional float64 VRD2
	optional float64 VOD2
	optional boolean Overvoltage2
	optional boolean Power2
	optional float64 VGATE3
	optional float64 VSUBS3
	optional float64 VRD3
	optional float64 VOD3
	optional boolean Overvoltage3
	optional boolean Power3

	optional int64  STATTime (TIMESTAMP(NANOS, true))
	optional int64  STATNanoseconds
	optional uint16 SPID
	optional uint8  SPREV
	optional uint16 FPID
	optional uint8  FPREV
	optional uint8  SVNA
	optional uint8  SVNB
	optional uint8  SVNC
	optional uint8  MODE
	optional uint32 EDACE
	optional uint32 EDACCE
	optional uint32 EDACN
	optional uint32 SPWEOP
	optional uint32 SPWEEP
	optional uint8  ANOMALY

	optional string TCV
	optional uint16 TCPID
	optional uint16 PSC
	optional uint8  ErrorCode

	optional group Warnings (LIST) {
		repeated group list {
			required string warning
		}
	}
	optional group Errors (LIST) {
		repeated group list {
			required string error
		}
	}
}`
