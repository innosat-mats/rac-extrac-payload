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

//Packages reads Ramses packages from buffer
func Packages(buf io.Reader, packageChannel chan<- Package, errorChannel chan<- error) {
	defer close(packageChannel)
	defer close(errorChannel)
	var err error
	for {
		header := Ramses{}
		err = header.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			errorChannel <- err
			break
		}

		if !header.Valid() {
			errorChannel <- errors.New("Not a valid RAC-file")
			break
		}

		secureHeader := Secure{}
		if header.SecureTrans() {
			err = secureHeader.Read(buf)
			if err != nil {
				errorChannel <- err
				break
			}
			header.Length -= uint16(binary.Size(secureHeader))
		}

		payload := make([]byte, header.Length)
		_, err = buf.Read(payload)
		if err != nil {
			errorChannel <- err
			break
		}
		if err == nil {
			packageChannel <- Package{header, secureHeader, payload}
		}
	}
}
