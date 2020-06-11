package main

func info() {
	println(`
### CSVs ###

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
  Counter of the transfer fram the payload packet arrived in.
- SPSequenceCount (Innosat Source Header)
- TMHeaderTime (Innosat TM Header)
  The time of the TM packet creation (UTC)
- TMHeaderNanoseconds
  The time of the TM packet creation (nanoseconds since epoch)
- SID
  The name of the SID or empty if the packet has no SID
- RID
  The name of the RID or empty if the packet has no RID

Note that each line should have a SID or a RID depending on packet type

All files also has a final common column
- Error
  If an error occurred it will be written here.
  If empty, then no error occurred extracting the data.

### CCD.csv ###
The following columns directly export the data in the rac:
CCDSEL, WDWOV, JPEGQ, FRAME, NROW, NRBIN, NRSKIP, NCOL, NCSKIP, NFLUSH
TEXPMS, TEMP, FBINOV, LBLNK, TBLNK, ZERO, TIMING1, TIMING2, VERSION
TIMING3, NBC, BC

The following columns parses the values further:
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
  The actual number of FPGA Columns (value in rac is the exponent in 2^x)
- NCBIN CCDColumns
  The number of CCD Columns
- GAIN Mode
  "High" (value in rac 0b0)
  "Low" (value in rac 0b1)
- GAIN Timing
  "Faster" (value in rac 0b0)
  "Full" (value in rac 0b1)
- GAIN Trunctation
  The value of the truncation bits
- Image File Name
  The name of the image file associated with these measurements
	`)
}
