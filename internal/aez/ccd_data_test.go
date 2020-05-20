package aez

import (
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
		EXPTS  uint32
		EXPTSS uint16
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
				EXPTS:  tt.fields.EXPTS,
				EXPTSS: tt.fields.EXPTSS,
			}
			if got := ccdImagePackData.Time(tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CCDImagePackData.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_Nanoseconds(t *testing.T) {
	type fields struct {
		EXPTS  uint32
		EXPTSS uint16
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
				EXPTS:  tt.fields.EXPTS,
				EXPTSS: tt.fields.EXPTSS,
			}
			if got := ccdImagePackData.Nanoseconds(); got != tt.want {
				t.Errorf("CCDImagePackData.Nanoseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_CSVHeaders_EqualLengthAs_CSVRow(t *testing.T) {
	ccd := CCDImagePackData{}
	headers := ccd.CSVHeaders()
	row := ccd.CSVRow()
	if len(headers) != len(row) {
		t.Errorf(
			"CCDImagePackData.CSVHeaders() length %v != CCDImagePackData.CSVRow() length %v",
			len(headers),
			len(row),
		)
	}
}

func TestWDWMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode WDWMode
		want string
	}{
		{"WDWModeAutomatic => Automatic", WDWModeAutomatic, "Automatic"},
		{"WDWModeManual => Manual", WDWModeManual, "Manual"},
		{"Else unrecognized", WDWMode(42), "Unrecognized WDWMode: 42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("WDWMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDGainMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode CCDGainMode
		want string
	}{
		{"HighSignalMode => High", HighSignalMode, "High"},
		{"LowSignalMode => Low", LowSignalMode, "Low"},
		{"Else unrecognized", CCDGainMode(42), "Unrecognized CCDGainMode: 42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("CCDGainMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCCDImagePackData_CSVRow(t *testing.T) {
	ccd := CCDImagePackData{
		CCDSEL:  5,
		EXPTS:   10,
		EXPTSS:  0xc000,
		WDW:     0x83,
		WDWOV:   13,
		JPEGQ:   101,
		FRAME:   14,
		NROW:    15,
		NRBIN:   16,
		NRSKIP:  17,
		NCOL:    18,
		NCBIN:   0xf648, // 0x6 -> 2^6 and then 0x48 -> 72
		NCSKIP:  19,
		NFLUSH:  20,
		TEXPMS:  21,
		GAIN:    0x1100, // Low and Full
		TEMP:    22,
		FBINOV:  23,
		LBLNK:   24,
		TBLNK:   25,
		ZERO:    26,
		TIMING1: 27,
		TIMING2: 28,
		VERSION: 29,
		TIMING3: 30,
		NBC:     31,
	}
	want := []string{
		"5",
		"10750000000",
		"1980-01-06T00:00:10.75Z",
		"Automatic",
		"14..3",
		"13",
		"101",
		"14",
		"15",
		"16",
		"17",
		"18",
		"64",
		"72",
		"19",
		"20",
		"21",
		"Low",
		"Full",
		"22",
		"23",
		"24",
		"25",
		"26",
		"27",
		"28",
		"29",
		"30",
		"31",
	}
	if got := ccd.CSVRow(); !reflect.DeepEqual(got, want) {
		t.Errorf("CCDImagePackData.CSVRow() = %v, want %v", got, want)
	}
}
