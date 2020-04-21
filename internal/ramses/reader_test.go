package ramses

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestPackages(t *testing.T) {
	type args struct {
		buf  io.Reader
		want Package
		err  error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"all zeros, not a vaild file",
			args{
				bytes.NewReader(make([]byte, 16)),
				Package{},
				errors.New("Not a valid RAC-file"),
			},
		},
		{
			"too short buffer",
			args{
				bytes.NewReader([]byte{0x90, 0xeb, 0}),
				Package{},
				errors.New("unexpected EOF"),
			},
		},
		{
			"a correct package 0 length no payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 16, 0}, make([]byte, 28+1)...)),
				Package{Payload: []byte{}},
				nil,
			},
		},
		{
			"a correct package 1 byte payload",
			args{
				bytes.NewReader(
					append([]byte{0x90, 0xeb, 17, 0}, make([]byte, 29+1)...)),
				Package{Payload: []byte{0}},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packs := make(chan Package)
			errors := make(chan error)
			go Packages(tt.args.buf, packs, errors)
			select {
			case got := <-packs:
				if !reflect.DeepEqual(got.Payload, tt.args.want.Payload) {
					t.Errorf("Package.Payload = %v, want %v", got.Payload, tt.args.want.Payload)
				}
			case err := <-errors:
				if err != nil {
					if tt.args.err == nil {
						t.Errorf("Package.Error = %v, want nil", err)
					} else {
						if err.Error() != tt.args.err.Error() {
							t.Errorf("Package.Error = %v, want %v", err, tt.args.err)
						}
					}
				}
			}

		})
	}
}
