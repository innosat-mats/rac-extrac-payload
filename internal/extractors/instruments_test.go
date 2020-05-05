package extractors

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_instrumentHK(t *testing.T) {
	type args struct {
		sid aez.SID
		buf io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    common.Exportable
		wantErr bool
	}{
		{
			"STAT",
			args{sid: aez.SIDSTAT, buf: bytes.NewReader(make([]byte, 100))},
			aez.STAT{},
			false,
		},
		{
			"HTR",
			args{sid: aez.SIDHTR, buf: bytes.NewReader(make([]byte, 100))},
			aez.HTR{},
			false,
		},
		{
			"PWR",
			args{sid: aez.SIDPWR, buf: bytes.NewReader(make([]byte, 100))},
			aez.PWR{},
			false,
		},
		{
			"CPRUA",
			args{sid: aez.SIDCPRUA, buf: bytes.NewReader(make([]byte, 100))},
			aez.CPRU{},
			false,
		},
		{
			"CPRUB",
			args{sid: aez.SIDCPRUB, buf: bytes.NewReader(make([]byte, 100))},
			aez.CPRU{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instrumentHK(tt.args.sid, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("instrumentHK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instrumentHK() = %v, want %v", got, tt.want)
			}
		})
	}
}
