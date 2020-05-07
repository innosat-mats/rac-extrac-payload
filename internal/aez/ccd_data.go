package aez

import (
	"encoding/binary"
	"fmt"
	"io"
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

// Wdw is the composite status of Image Window Mode
type Wdw uint8

// Mode returns the WDWMode type used encoded in Bit[7]
func (wdw Wdw) Mode() WDWMode {
	if (wdw & 0x80) == 0 {
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
func (wdw Wdw) InputDataWindow() (int, int, error) {
	switch wdw & 0b111 {
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
		"WDW value has unknown Input Data Window '%x'", wdw&0b111,
	)
}

// NCBin contains the FPGA and CCD columns bin count
type NCBin uint16

// FPGAColumns returns number FPGA columns to bin, Bit[11..8]
func (ncBin NCBin) FPGAColumns() int {
	return 1 << ((ncBin >> 8) & 0x0f)
}

// CCDColumns returns number of CCD columns to bin, Bit[7..0]
func (ncBin NCBin) CCDColumns() int {
	return (int)(ncBin & 0xff)
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

// CCDGainTiming is the timing flag
type CCDGainTiming int

const (
	// FasterTiming used for binned and discarded pixels
	FasterTiming CCDGainTiming = iota
	// FullTiming used even for pixels that are not read out
	FullTiming
)

// Mode returns high/low signal mode, Bit[12]
func (gain CCDGain) Mode() CCDGainMode {
	if (gain & 0x1000) == 0 {
		return HighSignalMode
	}
	return LowSignalMode
}

// Timing returns the full timing flag, Bit[8]
func (gain CCDGain) Timing() CCDGainTiming {
	if (gain & 0x100) == 0 {
		return FasterTiming
	}
	return FullTiming
}

// Truncation returns number of bits to be truncated (digital gain), Bit[3..0]
func (gain CCDGain) Truncation() uint8 {
	return uint8(gain & 0b1111)
}

// JPEGQUncompressed16bit is the value for non-12bit image data
var JPEGQUncompressed16bit = uint8(101)

// NCOLStartOffset says how many more columns than reported the actual columns are
var NCOLStartOffset uint16 = 1

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

// Read CCDImagePackData from the buffer
//
// returns the BC (bad columns) array and the error status.
func (ccdImagePackData *CCDImagePackData) Read(buf io.Reader) ([]uint16, error) {
	err := binary.Read(buf, binary.LittleEndian, ccdImagePackData)
	if err != nil {
		return nil, err
	}
	badColumns := make([]uint16, ccdImagePackData.NBC)
	return badColumns, binary.Read(buf, binary.LittleEndian, &badColumns)
}

// Time returns the measurement time in UTC
func (ccdImagePackData *CCDImagePackData) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(ccdImagePackData.EXPTS, ccdImagePackData.EXPTSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (ccdImagePackData *CCDImagePackData) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(ccdImagePackData.EXPTS, ccdImagePackData.EXPTSS)
}
