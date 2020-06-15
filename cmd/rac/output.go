package main

func infoGeneral() {
	println(`
### All CSVs ###

The first row:
- "CODE": Means the following column says what version of the code produced
  the output.
- "v0.2.2": The code version
- "RAMES": Means the following column says what RAMSES specification was used.
- "SPU045-S2:6F": The RAMSES version
- "INNOSAT": Means the following column says what INNOSAT specification was
  used.
- "IS-OSE-ICD-0005:1": The INNOSAT version
- "AEZ": Means the following column says what AEZ specification was used.
- "AEZICD002:H": The AEZ version

The header row starts with a couple of columns common to all output and then
follows columns specific to each file


- File
  The full path to the rac-file on the computer that produced the csv
- ProcessingDate
  The local time when the file was processed
- RamsesTime (Ramses Header)
  The time when the ramses file was created (UTC)

- QualityIndicator (Ramses TM Header)
  Indicates whether the transported data is complete or partial
  (0 = Complete, 1 = partial).
- LossFlag (Ramses TM Header)
  Used to indicate that a sequence discontinuity has been detected
- VCFrameCounter (Ramses TM Header)
  Counter of the transfer frame the payload packet arrived in.
  Wraps at (2^16)-1


- SPSequenceCount (Innosat Source Header)
  A counter that increases with each packet, may never short cycle and should
  wrap around to zero after 2^14-1


- TMHeaderTime (Innosat TM Header)
  The time of the TM packet creation (UTC)
- TMHeaderNanoseconds
  The time of the TM packet creation (nanoseconds since epoch)
- SID
  The name of the SID or empty if the packet has no SID
- RID
  The name of the RID or empty if the packet has no RID

Note that each line should have a SID or a RID depending on packet type

All files also have a final common column

- Error
  If an error occurred it will be written here.
  If empty, then no error occurred extracting the data.

For information about fields specific to a certain csv use any of these:

-help CCD, -help CPRU, -help HTR, -help PWR, -help STAT, -help TCV,
-help PM
	`)
}

func infoCCD() {
	println(`
### CCD.csv ###

The following columns directly export the data in the rac:
CCDSEL, WDWOV, JPEGQ, FRAME, NROW, NRBIN, NRSKIP, NCOL, NCSKIP, NFLUSH
TEXPMS, TEMP, FBINOV, LBLNK, TBLNK, ZERO, TIMING1, TIMING2, VERSION
TIMING3, NBC, BC

The following columns parse the values further:
- EXP Nanoseconds
  Time of exposure (nanoseconds since epoch)
- EXP Date
  Time of exposure (UTC)
- WDW Mode
  "Manual" (value in rac 0b0)
  "Automatic" (value in rac 0b1)
- WDW InputDataWindow
  Written as the from - to bits used in the original image
  "11..0" (value in rac 0x0)
  "12..1" (value in rac 0x1)
  "13..2" (value in rac 0x2)
  "14..3" (value in rac 0x3)
  "15..4" (value in rac 0x4)
  "15..0" which is the full image (value in rac 0x7)
- NCBIN FPGAColumns
  The actual number of FPGA Columns (value in rac is the exponent in 2^N)
- NCBIN CCDColumns
  The number of CCD Columns
- GAIN Mode
  "High" (value in rac 0b0)
  "Low" (value in rac 0b1)
- GAIN Timing
  "Faster" used for binned and discarded (value in rac 0b0)
  "Full" used even for pixels that are not read out (value in rac 0b1)
- GAIN Trunctation
  The value of the truncation bits
- Image File Name
  The name of the image file associated with these measurements
	`)
}

func infoCPRU() {
	println(`
### CPRU.csv ###
All voltages are the calculated float values of their respective type
according to the specification and not the raw encoded integer of the rac.

- VGATE0:       voltage
- VSUBS0:       voltage
- VRD0:         voltage
- VOD0:         voltage
- Overvoltage0: If over voltage fault registered (bool)
- Power0:       If power is enabled (bool)
- VGATE1:       voltage
- VSUBS1:       voltage
- VRD1:         voltage
- VOD1:         voltage
- Overvoltage1: If over voltage fault registered (bool)
- Power1:       If power is enabled (bool)
- VGATE2:       voltage
- VSUBS2:       voltage
- VRD2:         voltage
- VOD2:         voltage
- Overvoltage2: If over voltage fault registered (bool)
- Power2:       If power is enabled (bool)
- VGATE3:       voltage
- VSUBS3:       voltage
- VRD3:         voltage
- VOD3:         voltage
- Overvoltage3: If over voltage fault registered (bool)
- Power3:       If power is enabled (bool)
  `)
}

func infoHTR() {
	println(`
### HTR.csv ###
All voltages are the calculated float values of their respective type
according to the specification and not the raw encoded integer of the rac.

All temperatures are calculated from the specification and given in degrees
Celcius.

- HTR1A:    temperature,
- HTR1B:    temperature,
- HTR1OD:   voltage
- HTR2A:    temperature,
- HTR2B:    temperature,
- HTR2OD:   voltage
- HTR7A:    temperature,
- HTR7B:    temperature,
- HTR7OD:   voltage,
- HTR8A:    temperature,
- HTR8B:    temperature,
- HTR8OD:   voltage,
- WARNINGS: A summary of warnings from the temperature calculations.
  The warnings come from the interpolator and probably indicate the measured
  resistance is out of range.
  Each warning is separated by a '|' character.
  `)
}

func infoPWR() {
	println(`
### PWR.csv ###
All voltages are the calculated float values of their respective type
according to the specification and not the raw encoded integer of the rac.

All currents are calulated from the specification.

All temperatures are calculated from the specification and given in degrees
Celcius.

- PWRT:     temperature,
- PWRP32V:  voltage,
- PWRP32C:  current,
- PWRP16V:  voltage,
- PWRP16C:  current,
- PWRM16V:  voltage,
- PWRM16C:  current,
- PWRP3V3:  voltage,
- PWRP3C3:  current,
- WARNINGS: A summary of warnings from the temperature calculations.
  The warnings come from the interpolator and probably indicate the measured
  resistance is out of range.
  Each warning is separated by a '|' character.
  `)
}

func infoSTAT() {
	println(`
### STAT.csv ###

The following fields are read out exactly as they are encoded in the rac:
SPID, SPREV, FPID, FPREV, SVNA, SVNB, SVNC, MODE, EDACE, EDACCE, EDACN,
SPWEOP, SPWEEP, ANOMALY

The fields TS and TSS are replaced by:

- STATTIME: The time of the packet (UTC)
- STATNANO: The time of the packet (nanoseconds since epoch)
  `)
}

func infoTCV() {
	println(`
### TCV.csv ###

This contains all the four telecommand verification types

- TCV
  "Accept" for both accept success and fail
  "Exec" for both execute success and fail
- TCPID
  A copy of the Packet ID header field of the TC header
- PSC
  A copy of the Sequence Control Header field of the TC header
- ErrorCode
  Empty if success else the fail code
  `)
}

func infoPM() {
	println(`
### PM.csv ###

The following fields are read out exactly as they are encoded in the rac:
PM1A, PM1ACNTR, PM1B, PM1BCNTR, PM1S, PM1SCNTR, PM2A, PM2ACNTR, PM2B,
PM2BCNTR, PM2S, PM2SCNTR

The fields EXPTS and EXPTSS are replaced by:
- PMTIME: The exposure time (UTC)
- PMNANO: The exposure time (nanoseconds since epoch)
  `)
}

func infoSpace() {
	println(`
 +--------------------------------------------------------------------------------+
 |..                 .                                             ..    .        |
 |    .                              .                                   ..       |
 |                                        .   .             ..            .       |
 |..                                                         .        . ~.        |
 | .                            ..    .             __.--´|` + "`" + `--.__         ..      |
 | .                                          __.--´|__.--|--._ |` + "`" + `--.__     .     |
 |                                      __.--´|__.--|--.__|__.--|--.__|` + "`" + `--.      .|
 |                                __.--´|__.--|--.__|__.--|--.__|__.--|_.-´   .   |
 |                          __.--´|__.--|--.__|__.--|--.__|__.--|_.--´            |
 |                         -.__.--|--.__|__.--´--.__|__.--|_.,-'*\                |
 |                         |   ` + "`" + `--|__.--|--.__|__.--|_.--´  |*   |                |
 |..                       |         ` + "`" + `--|__.--|_.--´    |   |.---´                |
 |o+o++:~~~..              |     __        '-´*         |   |*   ` + "`" + `--.__           |
 |~::::++++++::~~~...      |    /  ` + "`" + `--.     |#          |   .` + "`" + `-.__     ` + "`" + `.         |
 |       ..~~~~~~:~~~~~~...|    |      \    |#    _.-\  |    |    ` + "`" + `--._/          |
 |          ..   ....~~~~~~|    \      |#   |#   /   /* |    |` + "`" + `-._.-_.-\          |
 | .         .             |     ` + "`" + `--.__/*   |#   \_.-*  |    |_.-´ /*  /          |
 | ...                      ` + "`" + `--.__    *     |#          |_   |_--._\_.-           |
 |                         ....   ` + "`" + `--.__    |*     __.--' ` + "`" + `--´                    |
 |~~~~~~...       .   . ..~::~~~~~      ` + "`" + `--.|__.--´~~~..                          |
 |..~+++++:~:~~.~~~~~~~~~.~~:~:+::~::~+:~~~~         ......                       |
 |..  .~~:~~::~~~:~~~::~:~::~~:+:~~.~~+:::~~..              . .                   |
 |+:~~~~~~~~.~~~~~~~:~:::++:~++::~..~~..~~~.~~.. ..             ~ .               |
 |~~.~~:~~.... . .  ..~..~~~~~~~~~.  ..      .   ...                 .            |
 |:~:~:+:~~  .   .. ..               ..                                   M A T S |
 +--------------------------------------------------------------------------------+
  `)
}
