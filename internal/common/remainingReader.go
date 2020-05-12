package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// RemainingReader unifies an underlying file pointer and a bytes.Reader
// so that both can know how much remains to read.
// Inspired by https://stackoverflow.com/a/41215144
type RemainingReader struct {
	reader *io.Reader
	read   int64
	size   int64
}

// Len returns the unread remaining portion of the data
func (reader *RemainingReader) Len() int64 {
	return reader.size - reader.read
}

// Read implements the io.Reader interface and remembers how much has been read.
func (reader *RemainingReader) Read(buf []byte) (int, error) {
	n, err := (*reader.reader).Read(buf)
	atomic.AddInt64(&reader.read, int64(n))
	return n, err
}

// NewRemainingReader returns a pointer to a RemaingReader.
// Only supports files and bytes.Reader
func NewRemainingReader(reader io.Reader) (*RemainingReader, error) {
	var size int64 = -1
	switch (reader).(type) {
	case *os.File:
		f, ok := reader.(*os.File)
		if !ok {
			return nil, errors.New("could not cast file")
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("could not stat file: %v", err)
		}
		size = fi.Size()
	case *bytes.Reader:
		r, ok := reader.(*bytes.Reader)
		if !ok {
			return nil, errors.New("could not cast bytes stream")
		}
		size = int64(r.Len())
	default:
		return nil, fmt.Errorf("no support creating a RemainingReader for %T", reader)
	}
	return &RemainingReader{reader: &reader, read: 0, size: size}, nil
}
