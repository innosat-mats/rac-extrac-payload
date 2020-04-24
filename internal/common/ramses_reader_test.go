package common

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDecodeRamses(t *testing.T) {
	type args struct {
		buf io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"all zeros, not a vaild file",
			args{bytes.NewReader(make([]byte, 16))},
			[]byte{},
			true,
		},
		{
			"too short buffer",
			args{bytes.NewReader([]byte{0x90, 0xeb, 0})},
			[]byte{},
			true,
		},
		{
			"a correct package 0 length no payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 16, 0}, make([]byte, 28+1)...))},
			[]byte{},
			false,
		},
		{
			"a correct package 1 byte payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 17, 0}, make([]byte, 29+1)...))},
			[]byte{0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packs := make(chan DataRecord)
			go DecodeRamses(packs, StreamBatch{Buf: tt.args.buf})
			got := <-packs
			if !reflect.DeepEqual(got.Buffer, tt.want) {
				t.Errorf("Package.Payload = %v, want %v", got.Buffer, tt.want)
			}
			if (got.Error != nil) != tt.wantErr {
				t.Errorf("DecodeSource() error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}

		})
	}
}
