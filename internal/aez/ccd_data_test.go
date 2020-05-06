package aez

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

func TestWdw_Mode(t *testing.T) {
	tests := []struct {
		name string
		wdw  Wdw
		want WDWMode
	}{
		{"Reads the correct bit", 0b10000000, WDWModeAutomatic},
		{"Is manual if bit is zero", 0, WDWModeManual},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wdw.Mode(); got != tt.want {
				t.Errorf("Wdw.Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWdw_InputDataWindow(t *testing.T) {
	tests := []struct {
		name    string
		wdw     Wdw
		want    int
		want1   int
		wantErr bool
	}{
		{"Returns expected 0x0 values", 0, 11, 0, false},
		{"Masks everything but Bit[2..0]", 0b111000, 11, 0, false},
		{"Returns expected 0x1 values", 0x1, 12, 1, false},
		{"Returns expected 0x2 values", 0x2, 13, 2, false},
		{"Returns expected 0x3 values", 0x3, 14, 3, false},
		{"Returns expected 0x4 values", 0x4, 15, 4, false},
		{"Returns all 16 bits for 0x7", 0x7, 15, 0, false},
		{"Returns error for other", 0x5, -1, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.wdw.InputDataWindow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Wdw.InputDataWindow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Wdw.InputDataWindow() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Wdw.InputDataWindow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNCBin_FPGAColumns(t *testing.T) {
	tests := []struct {
		name  string
		ncBin NCBin
		want  int
	}{
		{"Returns expected maximum 2^15", 0xffff, 32768},
		{"Returns expected minimum 2^0", 0, 1},
		{"Returns 2^3", 0b1100000000, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ncBin.FPGAColumns(); got != tt.want {
				t.Errorf("NCBin.FPGAColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNCBin_CCDColumns(t *testing.T) {
	tests := []struct {
		name  string
		ncBin NCBin
		want  int
	}{
		{"0x0000 returns 0", 0, 0},
		{"0xffff returns 255", 0xffff, 255},
		{"0x00ff returns 255", 0xff, 255},
		{"0xff00 returns 0", 0xff00, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ncBin.CCDColumns(); got != tt.want {
				t.Errorf("NCBin.CCDColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDGain_Mode(t *testing.T) {
	tests := []struct {
		name string
		gain CCDGain
		want CCDGainMode
	}{
		{"Reads 12th bit", 1 << 12, LowSignalMode},
		{"bit as 0 means high", 0, HighSignalMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gain.Mode(); got != tt.want {
				t.Errorf("CCDGain.Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDGain_Timing(t *testing.T) {
	tests := []struct {
		name string
		gain CCDGain
		want CCDGainTiming
	}{
		{"Reads the 8th bit", 1 << 8, FullTiming},
		{"Bit as 0 means faster timing", 0, FasterTiming},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gain.Timing(); got != tt.want {
				t.Errorf("CCDGain.Timing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDGain_Truncation(t *testing.T) {
	tests := []struct {
		name string
		gain CCDGain
		want uint8
	}{
		{"Reads Bit[3..0]", 0xabcd, 0xd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gain.Truncation(); got != tt.want {
				t.Errorf("CCDGain.Truncation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_Time(t *testing.T) {
	type fields struct {
		CCDSEL  uint8
		EXPTS   uint32
		EXPTSS  uint16
		WDW     Wdw
		WDWOV   uint16
		JPEGQ   uint8
		FRAME   uint16
		NROW    uint16
		NRBIN   uint16
		NRSKIP  uint16
		NCOL    uint16
		NCBIN   NCBin
		NCSKIP  uint16
		NFLUSH  uint16
		TEXPMS  uint32
		GAIN    CCDGain
		TEMP    uint16
		FBINOV  uint16
		LBLNK   uint16
		TBLNK   uint16
		ZERO    uint16
		TIMING1 uint16
		TIMING2 uint16
		VERSION uint16
		TIMING3 uint16
		NBC     uint16
	}
	type args struct {
		epoch time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			"Generates expected time",
			fields{EXPTS: 10, EXPTSS: 0b1100000000000000},
			args{epoch: ccsds.TAI},
			ccsds.TAI.Add(time.Second * 10).Add(time.Millisecond * 750),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccdImagePackData := &CCDImagePackData{
				CCDSEL:  tt.fields.CCDSEL,
				EXPTS:   tt.fields.EXPTS,
				EXPTSS:  tt.fields.EXPTSS,
				WDW:     tt.fields.WDW,
				WDWOV:   tt.fields.WDWOV,
				JPEGQ:   tt.fields.JPEGQ,
				FRAME:   tt.fields.FRAME,
				NROW:    tt.fields.NROW,
				NRBIN:   tt.fields.NRBIN,
				NRSKIP:  tt.fields.NRSKIP,
				NCOL:    tt.fields.NCOL,
				NCBIN:   tt.fields.NCBIN,
				NCSKIP:  tt.fields.NCSKIP,
				NFLUSH:  tt.fields.NFLUSH,
				TEXPMS:  tt.fields.TEXPMS,
				GAIN:    tt.fields.GAIN,
				TEMP:    tt.fields.TEMP,
				FBINOV:  tt.fields.FBINOV,
				LBLNK:   tt.fields.LBLNK,
				TBLNK:   tt.fields.TBLNK,
				ZERO:    tt.fields.ZERO,
				TIMING1: tt.fields.TIMING1,
				TIMING2: tt.fields.TIMING2,
				VERSION: tt.fields.VERSION,
				TIMING3: tt.fields.TIMING3,
				NBC:     tt.fields.NBC,
			}
			if got := ccdImagePackData.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImagePackData.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_Nanoseconds(t *testing.T) {
	type fields struct {
		CCDSEL  uint8
		EXPTS   uint32
		EXPTSS  uint16
		WDW     Wdw
		WDWOV   uint16
		JPEGQ   uint8
		FRAME   uint16
		NROW    uint16
		NRBIN   uint16
		NRSKIP  uint16
		NCOL    uint16
		NCBIN   NCBin
		NCSKIP  uint16
		NFLUSH  uint16
		TEXPMS  uint32
		GAIN    CCDGain
		TEMP    uint16
		FBINOV  uint16
		LBLNK   uint16
		TBLNK   uint16
		ZERO    uint16
		TIMING1 uint16
		TIMING2 uint16
		VERSION uint16
		TIMING3 uint16
		NBC     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"Generates expected nanoseconds",
			fields{EXPTS: 10, EXPTSS: 0b1100000000000000},
			10750000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccdImagePackData := &CCDImagePackData{
				CCDSEL:  tt.fields.CCDSEL,
				EXPTS:   tt.fields.EXPTS,
				EXPTSS:  tt.fields.EXPTSS,
				WDW:     tt.fields.WDW,
				WDWOV:   tt.fields.WDWOV,
				JPEGQ:   tt.fields.JPEGQ,
				FRAME:   tt.fields.FRAME,
				NROW:    tt.fields.NROW,
				NRBIN:   tt.fields.NRBIN,
				NRSKIP:  tt.fields.NRSKIP,
				NCOL:    tt.fields.NCOL,
				NCBIN:   tt.fields.NCBIN,
				NCSKIP:  tt.fields.NCSKIP,
				NFLUSH:  tt.fields.NFLUSH,
				TEXPMS:  tt.fields.TEXPMS,
				GAIN:    tt.fields.GAIN,
				TEMP:    tt.fields.TEMP,
				FBINOV:  tt.fields.FBINOV,
				LBLNK:   tt.fields.LBLNK,
				TBLNK:   tt.fields.TBLNK,
				ZERO:    tt.fields.ZERO,
				TIMING1: tt.fields.TIMING1,
				TIMING2: tt.fields.TIMING2,
				VERSION: tt.fields.VERSION,
				TIMING3: tt.fields.TIMING3,
				NBC:     tt.fields.NBC,
			}
			if got := ccdImagePackData.Nanoseconds(); got != tt.want {
				t.Errorf("CCDImagePackData.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_Read(t *testing.T) {
	type fields struct {
		CCDSEL  uint8
		EXPTS   uint32
		EXPTSS  uint16
		WDW     Wdw
		WDWOV   uint16
		JPEGQ   uint8
		FRAME   uint16
		NROW    uint16
		NRBIN   uint16
		NRSKIP  uint16
		NCOL    uint16
		NCBIN   NCBin
		NCSKIP  uint16
		NFLUSH  uint16
		TEXPMS  uint32
		GAIN    CCDGain
		TEMP    uint16
		FBINOV  uint16
		LBLNK   uint16
		TBLNK   uint16
		ZERO    uint16
		TIMING1 uint16
		TIMING2 uint16
		VERSION uint16
		TIMING3 uint16
		NBC     uint16
	}
	tests := []struct {
		name      string
		fields    fields
		extraData []uint16
		want      []uint16
		wantErr   bool
	}{
		{
			"Errors if NBC is too large for remaining buffer",
			fields{NBC: 10},
			[]uint16{1, 2, 3, 4, 5},
			[]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			true,
		},
		{
			"Returns an array of bad columns data",
			fields{NBC: 3},
			[]uint16{1, 2, 3, 4, 5},
			[]uint16{1, 2, 3},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccdImagePackData := &CCDImagePackData{
				CCDSEL:  tt.fields.CCDSEL,
				EXPTS:   tt.fields.EXPTS,
				EXPTSS:  tt.fields.EXPTSS,
				WDW:     tt.fields.WDW,
				WDWOV:   tt.fields.WDWOV,
				JPEGQ:   tt.fields.JPEGQ,
				FRAME:   tt.fields.FRAME,
				NROW:    tt.fields.NROW,
				NRBIN:   tt.fields.NRBIN,
				NRSKIP:  tt.fields.NRSKIP,
				NCOL:    tt.fields.NCOL,
				NCBIN:   tt.fields.NCBIN,
				NCSKIP:  tt.fields.NCSKIP,
				NFLUSH:  tt.fields.NFLUSH,
				TEXPMS:  tt.fields.TEXPMS,
				GAIN:    tt.fields.GAIN,
				TEMP:    tt.fields.TEMP,
				FBINOV:  tt.fields.FBINOV,
				LBLNK:   tt.fields.LBLNK,
				TBLNK:   tt.fields.TBLNK,
				ZERO:    tt.fields.ZERO,
				TIMING1: tt.fields.TIMING1,
				TIMING2: tt.fields.TIMING2,
				VERSION: tt.fields.VERSION,
				TIMING3: tt.fields.TIMING3,
				NBC:     tt.fields.NBC,
			}
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.fields)
			if err != nil {
				t.Errorf("Could not inject struct data into buffer")
				return
			}
			err = binary.Write(buf, binary.LittleEndian, tt.extraData)
			if err != nil {
				t.Errorf("Could not inject extra data into buffer")
				return
			}
			got, err := ccdImagePackData.Read(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("CCDImagePackData.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImagePackData.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
