package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

// WDWMode describes the CCD WDW parameter
type WDWMode uint8

const (
	// WDWModeManual check WDWOV for overflows
	WDWModeManual WDWMode = 0
	// WDWModeAutomatic window selected depending on input data
	WDWModeAutomatic WDWMode = 1
)

func (mode *WDWMode) String() string {
	switch *mode {
	case WDWModeManual:
		return "Manual"
	case WDWModeAutomatic:
		return "Automatic"
	default:
		return fmt.Sprintf("Unrecognized WDWMode: %v", int(*mode))
	}
}

// Wdw is the composite status of Image Window Mode
type Wdw uint8

// Mode returns the WDWMode type used encoded in Bit[7]
func (wdw *Wdw) Mode() WDWMode {
	if (*wdw & 0x80) == 0 {
		return WDWModeManual
	}
	return WDWModeAutomatic
}

// InputDataWindow returns which bits out of the depth of the CCD Images
//
// It ueses information in Bit[2..0]
//
// Returns major and minor bit of original CCD Image and error
// if encoding not covered by specification.
//
// If full range of Bit[15..0] is used the JPEGQ should be 0xFF
func (wdw *Wdw) InputDataWindow() (int, int, error) {
	switch *wdw & 0b111 {
	case 0x0:
		return 11, 0, nil
	case 0x1:
		return 12, 1, nil
	case 0x2:
		return 13, 2, nil
	case 0x3:
		return 14, 3, nil
	case 0x4:
		return 15, 4, nil
	case 0x7:
		return 15, 0, nil
	}
	return -1, -1, fmt.Errorf(
		"WDW value has unknown Input Data Window '%x'", *wdw&0b111,
	)
}

// NCBin contains the FPGA and CCD columns bin count
type NCBin uint16

// FPGAColumns returns number FPGA columns to bin, Bit[11..8]
//  the value is encoded as 2^x
func (ncBin *NCBin) FPGAColumns() int {
	return 1 << ((*ncBin >> 8) & 0x0f)
}

// CCDColumns returns number of CCD columns to bin, Bit[7..0]
func (ncBin *NCBin) CCDColumns() int {
	return (int)(*ncBin & 0xff)
}

// CCDGain is game composite information
type CCDGain uint16

// CCDGainMode is high/low signal mode
type CCDGainMode int

// HighSignalMode, LowSignalMode
const (
	HighSignalMode CCDGainMode = iota
	LowSignalMode
)

func (mode *CCDGainMode) String() string {
	switch *mode {
	case HighSignalMode:
		return "High"
	case LowSignalMode:
		return "Low"
	default:
		return fmt.Sprintf("Unrecognized CCDGainMode: %v", int(*mode))
	}
}

// CCDGainTiming is the timing flag
type CCDGainTiming int

const (
	// FasterTiming used for binned and discarded pixels
	FasterTiming CCDGainTiming = iota
	// FullTiming used even for pixels that are not read out
	FullTiming
)

func (timing *CCDGainTiming) String() string {
	switch *timing {
	case FasterTiming:
		return "Faster"
	case FullTiming:
		return "Full"
	default:
		return ""
	}
}

// Mode returns high/low signal mode, Bit[12]
func (gain *CCDGain) Mode() CCDGainMode {
	if (*gain & 0x1000) == 0 {
		return HighSignalMode
	}
	return LowSignalMode
}

// Timing returns the full timing flag, Bit[8]
func (gain *CCDGain) Timing() CCDGainTiming {
	if (*gain & 0x100) == 0 {
		return FasterTiming
	}
	return FullTiming
}

// Truncation returns number of bits to be truncated (digital gain), Bit[3..0]
func (gain *CCDGain) Truncation() uint8 {
	return uint8(*gain & 0b1111)
}

// JPEGQUncompressed16bit is the value for non-12bit image data
const JPEGQUncompressed16bit = uint8(101)

// NCOLStartOffset says how many more columns than reported the actual columns are
const NCOLStartOffset uint16 = 1

// CCDImagePackData contains the composite information from the CCD and the CRB module
type CCDImagePackData struct {
	CCDSEL  uint8   // CCD sensor number, not same format for TM and TC
	EXPTS   uint32  // Exposure start time, seconds (CUC time format)
	EXPTSS  uint16  // Exposure start time, subseconds (CUC time format)
	WDW     Wdw     // Window mode
	WDWOV   uint16  // Bit window overflow counter (should be zero when WDW is automatic)
	JPEGQ   uint8   // JPEG compression quality setting (0..100)
	FRAME   uint16  // Frame count since boot
	NROW    uint16  // Number of rows in image 1..511
	NRBIN   uint16  // Number of rows to bin together 0..63
	NRSKIP  uint16  // Number of rows to skip before start of readout 0..511
	NCOL    uint16  // Number of columns in image (starts at 0) 0..2047
	NCBIN   NCBin   // Number of columns to bin in FPGA (power of two) Bit[11..8], CCD Bit[7..0]
	NCSKIP  uint16  // Numbers of columns to skip before start of readout
	NFLUSH  uint16  // Number of pre-exposure flushes
	TEXPMS  uint32  // Exposure time in milliseconds
	GAIN    CCDGain // Gain composite information
	TEMP    uint16  // Temperature of the ADC
	FBINOV  uint16  // Number of overflows detected while binning 0..32767
	LBLNK   uint16  // Value of leading blanks. Average of blank pixesl 32:47 from the middle row of the readout region
	TBLNK   uint16  // Value of trailing blanks. Average of blank pixels 2128:2143 (50 leading blank + 2048 pxiels + 50 trailing blanks) from the middle row of the readout region
	ZERO    uint16  // Value of zero input reading
	TIMING1 uint16  // Clock timing parameters, Bit[15..0]. Alternates with frame count between *1rt and *2rt timings
	TIMING2 uint16  // Clock timing parameters, Bit[31..16]. Alternates with frame count between *1rt and *2rt timings
	VERSION uint16  // Readout of firmware version
	TIMING3 uint16  // Clock timing parameters, Bit[47..32]. Alternates with frame count between *1rt and *2rt timings
	NBC     uint16  // Number of bad columns set
}

// NewCCDImagePackData reads buf into CCDImagePackData
func NewCCDImagePackData(buf io.Reader) (*CCDImagePackData, error) {
	ccd := CCDImagePackData{}
	err := binary.Read(buf, binary.LittleEndian, &ccd)
	if err != nil {
		return nil, err
	}
	return &ccd, nil
}

// Time returns the measurement time in UTC
func (ccd *CCDImagePackData) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(ccd.EXPTS, ccd.EXPTSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (ccd *CCDImagePackData) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(ccd.EXPTS, ccd.EXPTSS)
}

// CSVHeaders returns the exportable field names
func (ccd *CCDImagePackData) CSVHeaders() []string {
	return []string{
		"CCDSEL",
		"EXP Nanoseconds",
		"EXP Date",
		"WDW Mode",
		"WDW InputDataWindow",
		"WDWOV",
		"JPEGQ",
		"FRAME",
		"NROW",
		"NRBIN",
		"NRSKIP",
		"NCOL",
		"NCBIN FPGAColumns",
		"NCBIN CCDColumns",
		"NCSKIP",
		"NFLUSH",
		"TEXPMS",
		"GAIN Mode",
		"GAIN Timing",
		"GAIN Trunctation",
		"TEMP",
		"FBINOV",
		"LBLNK",
		"TBLNK",
		"ZERO",
		"TIMING1",
		"TIMING2",
		"VERSION",
		"TIMING3",
		"NBC",
	}
}

// CSVRow returns the exportable field values
func (ccd *CCDImagePackData) CSVRow() []string {
	wdwhigh, wdwlow, _ := ccd.WDW.InputDataWindow()
	wdwMode := ccd.WDW.Mode()
	gainMode := ccd.GAIN.Mode()
	gainTiming := ccd.GAIN.Timing()
	return []string{
		strconv.Itoa(int(ccd.CCDSEL)),
		strconv.FormatInt(ccd.Nanoseconds(), 10),
		ccd.Time(gpsTime).Format(time.RFC3339Nano),
		(&wdwMode).String(),
		fmt.Sprintf("%v..%v", wdwhigh, wdwlow),
		strconv.Itoa(int(ccd.WDWOV)),
		strconv.Itoa(int(ccd.JPEGQ)),
		strconv.Itoa(int(ccd.FRAME)),
		strconv.Itoa(int(ccd.NROW)),
		strconv.Itoa(int(ccd.NRBIN)),
		strconv.Itoa(int(ccd.NRSKIP)),
		strconv.Itoa(int(ccd.NCOL)),
		strconv.Itoa(int(ccd.NCBIN.FPGAColumns())),
		strconv.Itoa(int(ccd.NCBIN.CCDColumns())),
		strconv.Itoa(int(ccd.NCSKIP)),
		strconv.Itoa(int(ccd.NFLUSH)),
		strconv.Itoa(int(ccd.TEXPMS)),
		(&gainMode).String(),
		(&gainTiming).String(),
		strconv.Itoa(int(ccd.GAIN.Truncation())),
		strconv.Itoa(int(ccd.TEMP)),
		strconv.Itoa(int(ccd.FBINOV)),
		strconv.Itoa(int(ccd.LBLNK)),
		strconv.Itoa(int(ccd.TBLNK)),
		strconv.Itoa(int(ccd.ZERO)),
		strconv.Itoa(int(ccd.TIMING1)),
		strconv.Itoa(int(ccd.TIMING2)),
		strconv.Itoa(int(ccd.VERSION)),
		strconv.Itoa(int(ccd.TIMING3)),
		strconv.Itoa(int(ccd.NBC)),
	}
}
