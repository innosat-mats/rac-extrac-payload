package ramses

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// Package is complete Ramses package
type Package struct {
	Header       Ramses
	SecureHeader Secure
	Payload      []byte
}

//Packages ...
func Packages(buf io.Reader, rPack chan Package, rErr chan error) {
	defer close(rPack)
	defer close(rErr)
	var err error
	for {
		header := Ramses{}
		err = header.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			rErr <- err
			break
		}

		if !header.Valid() {
			rErr <- errors.New("Not a valid RAC-file")
			break
		}

		secureHeader := Secure{}
		if header.SecureTrans() {
			err = secureHeader.Read(buf)
			if err != nil {
				rErr <- err
				break
			}
			header.Length -= uint16(binary.Size(secureHeader))
		}

		payload := make([]byte, header.Length)
		_, err = buf.Read(payload)
		if err != nil {
			rErr <- err
			break
		}
		if err == nil {
			rPack <- Package{header, secureHeader, payload}
		}
	}
}
