package timeseries

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

func parseTime(timestamp string) time.Time {
	layout := "2006-01-02 15:04:05 -0700 UTC"
	t, _ := time.Parse(layout, timestamp)
	return t.UTC()
}

func getTestImage() []byte {
	fileContents, err := os.ReadFile(filepath.Join("..", "aez", "testdata", "3166_4052_5.jpg"))
	if err != nil {
		log.Fatalln(err)
	}
	return fileContents
}

func TestGetParquetRow(t *testing.T) {
	procDate := time.Now()
	data := common.DataRecord{
		Origin:         &common.OriginDescription{Name: "Sputnik", ProcessingDate: procDate},
		RamsesHeader:   &ramses.Ramses{Date: 24, Time: 42000},
		RamsesTMHeader: &ramses.TMHeader{LossFlag: 1, VCFrameCounter: 42},
		SourceHeader:   &innosat.SourcePacketHeader{PacketSequenceControl: innosat.PacketSequenceControl(0xc003)},
		TMHeader:       &innosat.TMHeader{CUCTimeSeconds: 42, CUCTimeFraction: 0xc000},
		Buffer:         getTestImage(),
	}
	type args struct {
		sid  aez.SID
		data common.Exporter
	}
	tests := []struct {
		name string
		args args
		want parquetrow.ParquetRow
	}{
		{
			"Test CCDImage",
			args{0, &aez.CCDImage{
				PackData:      &aez.CCDImagePackData{JPEGQ: 95},
				ImageFileName: "HelloWorld.png",
			}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				EXPDate:             parseTime("1980-01-05 23:59:42 +0000 UTC"),
				WDWMode:             "Manual",
				WDWInputDataWindow:  "11..0",
				JPEGQ:               95,
				NCBINFPGAColumns:    1,
				GAINMode:            "High",
				GAINTiming:          "Faster",
				ImageName:           "HelloWorld.png",
			},
		},
		{
			"Test PMData",
			args{0, &aez.PMData{}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				PMTime:              parseTime("1980-01-05 23:59:42 +0000 UTC"),
			},
		},
		{
			"Test HTR",
			args{aez.SIDHTR, &aez.HTR{}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "HTR",
				RID:                 "",
				HTR1A:               -55,
				HTR1B:               -55,
				HTR2A:               -55,
				HTR2B:               -55,
				HTR7A:               -55,
				HTR7B:               -55,
				HTR8A:               -55,
				HTR8B:               -55,
				Warnings: []string{
					"HTR1A: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR1B: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR2A: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR2B: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR7A: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR7B: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR8A: +Inf is too large for interpolator. Returning value for maximum.",
					"HTR8B: +Inf is too large for interpolator. Returning value for maximum.",
				},
			},
		},
		{
			"Test PWR",
			args{aez.SIDPWR, &aez.PWR{}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "PWR",
				RID:                 "",
				PWRT:                -55,
				Warnings:            []string{"PWRT: +Inf is too large for interpolator. Returning value for maximum."},
			},
		},
		{
			"Test CPRU",
			args{aez.SIDCPRUA, &aez.CPRU{STAT: 8}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "CPRUA",
				RID:                 "",
				Power0:              true,
			},
		},
		{
			"Test STAT",
			args{aez.SIDSTAT, &aez.STAT{}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "STAT",
				RID:                 "",
				STATTime:            parseTime("1980-01-05 23:59:42 +0000 UTC"),
			},
		},
		{
			"Test TCAcceptSuccessData",
			args{0, &aez.TCAcceptSuccessData{TCPID: 1, PSC: 2}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				TCV:                 "Accept",
				TCPID:               1,
				PSC:                 2,
			},
		},
		{
			"Test TCAcceptFailureData",
			args{0, &aez.TCAcceptFailureData{TCPID: 1, PSC: 2, ErrorCode: 3}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				TCV:                 "Accept",
				TCPID:               1,
				PSC:                 2,
				ErrorCode:           3,
			},
		},
		{
			"Test TCExecSuccessData",
			args{0, &aez.TCExecSuccessData{TCPID: 1, PSC: 2}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				TCV:                 "Exec",
				TCPID:               1,
				PSC:                 2,
			},
		},
		{
			"Test TCExecFailureData",
			args{0, &aez.TCExecFailureData{TCPID: 1, PSC: 2, ErrorCode: 3}},
			parquetrow.ParquetRow{
				OriginFile:          "Sputnik",
				ProcessingTime:      procDate,
				RamsesTime:          data.RamsesHeader.Created(),
				QualityIndicator:    0,
				LossFlag:            1,
				VCFrameCounter:      42,
				SPSequenceCount:     3,
				TMHeaderTime:        data.TMHeader.Time(aez.GpsTime),
				TMHeaderNanoseconds: 42750000000,
				SID:                 "",
				RID:                 "",
				TCV:                 "Exec",
				TCPID:               1,
				PSC:                 2,
				ErrorCode:           3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.SID = tt.args.sid
			data.Data = tt.args.data
			if got := GetParquetRow(&data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParquetRow():\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}
