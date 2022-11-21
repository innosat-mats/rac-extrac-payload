package timeseries

// RacSchema is the parquet schema for saving RAC data, one row per packet
const RacSchema = `message schema {
	required binary OriginFile (STRING);
	required int64  ProcessingTime (TIMESTAMP(NANOS, true));
	required int64  RamsesTime (TIMESTAMP(NANOS, true));
	required int32  QualityIndicator;
	required int32  LossFlag;
	required int32  VCFrameCounter;
	required int32  SPSequenceCount;
	required int64  TMHeaderTime (TIMESTAMP(NANOS, true));
	required int64  TMHeaderNanoseconds;
	required binary SID (STRING);
	required binary RID (STRING);

	optional int32  CCDSEL;
	optional int64  EXPNanoseconds;
	optional int64  EXPDate (TIMESTAMP(NANOS, true));
	optional binary WDWMode (STRING);
	optional binary WDWInputDataWindow (STRING);
	optional int32  WDWOV;
	optional int32  JPEGQ;
	optional int32  FRAME;
	optional int32  NROW;
	optional int32  NRBIN;
	optional int32  NRSKIP;
	optional int32  NCOL;
	optional int32  NCBINFPGAColumns;
	optional int32  NCBINCCDColumns;
	optional int32  NCSKIP;
	optional int32  NFLUSH;
	optional int32  TEXPMS;
	optional binary GAINMode (STRING);
	optional binary GAINTiming (STRING);
	optional int32  GAINTruncation;
	optional int32  TEMP;
	optional int32  FBINOV;
	optional int32  LBLNK;
	optional int32  TBLNK;
	optional int32  ZERO;
	optional int32  TIMING1;
	optional int32  TIMING2;
	optional int32  VERSION;
	optional int32  TIMING3;
	optional int32  NBC;
	optional group  BadColumns (LIST) {
		repeated group list {
			required int32 element;
		}
	}
	optional binary ImageName (STRING);
	optional binary ImageData;

	optional int64 PMTime (TIMESTAMP(NANOS, true));
	optional int64 PMNanoseconds;
	optional int32 PM1A;
	optional int32 PM1ACNTR;
	optional int32 PM1B;
	optional int32 PM1BCNTR;
	optional int32 PM1S;
	optional int32 PM1SCNTR;
	optional int32 PM2A;
	optional int32 PM2ACNTR;
	optional int32 PM2B;
	optional int32 PM2BCNTR;
	optional int32 PM2S;
	optional int32 PM2SCNTR;

	optional double HTR1A;
	optional double HTR1B;
	optional double HTR1OD;
	optional double HTR2A;
	optional double HTR2B;
	optional double HTR2OD;
	optional double HTR7A;
	optional double HTR7B;
	optional double HTR7OD;
	optional double HTR8A;
	optional double HTR8B;
	optional double HTR8OD;

	optional double PWRT;
	optional double PWRP32V;
	optional double PWRP32C;
	optional double PWRP16V;
	optional double PWRP16C;
	optional double PWRM16V;
	optional double PWRM16C;
	optional double PWRP3V3;
	optional double PWRP3C3;

	optional double  VGATE0;
	optional double  VSUBS0;
	optional double  VRD0;
	optional double  VOD0;
	optional boolean Overvoltage0;
	optional boolean Power0;
	optional double  VGATE1;
	optional double  VSUBS1;
	optional double  VRD1;
	optional double  VOD1;
	optional boolean Overvoltage1;
	optional boolean Power1;
	optional double  VGATE2;
	optional double  VSUBS2;
	optional double  VRD2;
	optional double  VOD2;
	optional boolean Overvoltage2;
	optional boolean Power2;
	optional double  VGATE3;
	optional double  VSUBS3;
	optional double  VRD3;
	optional double  VOD3;
	optional boolean Overvoltage3;
	optional boolean Power3;

	optional int64 STATTime (TIMESTAMP(NANOS, true));
	optional int64 STATNanoseconds;
	optional int32 SPID;
	optional int32 SPREV;
	optional int32 FPID;
	optional int32 FPREV;
	optional int32 SVNA;
	optional int32 SVNB;
	optional int32 SVNC;
	optional int32 MODE;
	optional int32 EDACE;
	optional int32 EDACCE;
	optional int32 EDACN;
	optional int32 SPWEOP;
	optional int32 SPWEEP;
	optional int32 ANOMALY;

	optional binary TCV (STRING);
	optional int32  TCPID;
	optional int32  PSC;
	optional int32  ErrorCode;

	optional group Warnings (LIST) {
		repeated group list {
			required binary element (STRING);
		}
	}
	optional group Errors (LIST) {
		repeated group list {
			required binary element (STRING);
		}
	}
}`
