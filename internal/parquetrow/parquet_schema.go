package parquetrow

// RacCCDSchema is the parquet schema for saving RAC CCD data, one row per packet
const RacCCDSchema = `message schema {
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

	required int32  CCDSEL;
	required int64  EXPNanoseconds;
	required int64  EXPDate (TIMESTAMP(MICROS, true));
	required binary WDWMode (STRING);
	required binary WDWInputDataWindow (STRING);
	required int32  WDWOV;
	required int32  JPEGQ;
	required int32  FRAME;
	required int32  NROW;
	required int32  NRBIN;
	required int32  NRSKIP;
	required int32  NCOL;
	required int32  NCBINFPGAColumns;
	required int32  NCBINCCDColumns;
	required int32  NCSKIP;
	required int32  NFLUSH;
	required int32  TEXPMS;
	required binary GAINMode (STRING);
	required binary GAINTiming (STRING);
	required int32  GAINTruncation;
	required int32  TEMP;
	required int32  FBINOV;
	required int32  LBLNK;
	required int32  TBLNK;
	required int32  ZERO;
	required int32  TIMING1;
	required int32  TIMING2;
	required int32  VERSION;
	required int32  TIMING3;
	required int32  NBC;
	required group  BadColumns (LIST) {
		repeated group list {
			required int32 element;
		}
	}
	required binary ImageName (STRING);
	required binary ImageData;

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

// RacPMSchema is the parquet schema for saving RAC PM data, one row per packet
const RacPMSchema = `message schema {
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

	required int64 PMTime (TIMESTAMP(NANOS, true));
	required int64 PMNanoseconds;
	required int32 PM1A;
	required int32 PM1ACNTR;
	required int32 PM1B;
	required int32 PM1BCNTR;
	required int32 PM1S;
	required int32 PM1SCNTR;
	required int32 PM2A;
	required int32 PM2ACNTR;
	required int32 PM2B;
	required int32 PM2BCNTR;
	required int32 PM2S;
	required int32 PM2SCNTR;

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

// RacHTRSchema is the parquet schema for saving RAC HTR data, one row per packet
const RacHTRSchema = `message schema {
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

	required double HTR1A;
	required double HTR1B;
	required double HTR1OD;
	required double HTR2A;
	required double HTR2B;
	required double HTR2OD;
	required double HTR7A;
	required double HTR7B;
	required double HTR7OD;
	required double HTR8A;
	required double HTR8B;
	required double HTR8OD;

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

// RacPWRSchema is the parquet schema for saving RAC PWR data, one row per packet
const RacPWRSchema = `message schema {
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

	required double PWRT;
	required double PWRP32V;
	required double PWRP32C;
	required double PWRP16V;
	required double PWRP16C;
	required double PWRM16V;
	required double PWRM16C;
	required double PWRP3V3;
	required double PWRP3C3;

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

// RacCPRUSchema is the parquet schema for saving RAC CPRU data, one row per packet
const RacCPRUSchema = `message schema {
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

	required double  VGATE0;
	required double  VSUBS0;
	required double  VRD0;
	required double  VOD0;
	required boolean Overvoltage0;
	required boolean Power0;
	required double  VGATE1;
	required double  VSUBS1;
	required double  VRD1;
	required double  VOD1;
	required boolean Overvoltage1;
	required boolean Power1;
	required double  VGATE2;
	required double  VSUBS2;
	required double  VRD2;
	required double  VOD2;
	required boolean Overvoltage2;
	required boolean Power2;
	required double  VGATE3;
	required double  VSUBS3;
	required double  VRD3;
	required double  VOD3;
	required boolean Overvoltage3;
	required boolean Power3;

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

// RacSTATSchema is the parquet schema for saving RAC STAT data, one row per packet
const RacSTATSchema = `message schema {
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

	required int64 STATTime (TIMESTAMP(NANOS, true));
	required int64 STATNanoseconds;
	required int32 SPID;
	required int32 SPREV;
	required int32 FPID;
	required int32 FPREV;
	required int32 SVNA;
	required int32 SVNB;
	required int32 SVNC;
	required int32 MODE;
	required int32 EDACE;
	required int32 EDACCE;
	required int32 EDACN;
	required int32 SPWEOP;
	required int32 SPWEEP;
	required int32 ANOMALY;

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

// RacTCVSchema is the parquet schema for saving RAC TCV data, one row per packet
const RacTCVSchema = `message schema {
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

	required binary TCV (STRING);
	required int32  TCPID;
	required int32  PSC;
	required int32  ErrorCode;

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
	optional int64  EXPDate (TIMESTAMP(MICROS, true));
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
