package common

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewRemainingReader_CreatesFromFile(t *testing.T) {
	file, err := ioutil.TempFile("", "test-file")
	defer os.Remove(file.Name())
	file.Write([]byte("Hello!"))
	file.Seek(0, 0)
	reader, err := NewRemainingReader(file)
	if err != nil {
		t.Errorf("NewRemainingReader() returned unexpected error = %v", err)
	}
	if reader.Len() != 6 {
		t.Errorf("RemainingReader.Len() = %v, want %v", reader.Len(), 6)
	}
}

func TestNewRemainingReader_CreatesFromBytesReader(t *testing.T) {
	var buf = []byte("hello!")
	reader, err := NewRemainingReader(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("NewRemainingReader() returned unexpected error = %v", err)
	}
	if reader.Len() != 6 {
		t.Errorf("RemainingReader.Len() = %v, want %v", reader.Len(), 6)
	}
}

func TestNewRemainingReader_ErrorsAsDefault(t *testing.T) {
	var buf = []byte("hello!")
	_, err := NewRemainingReader(bytes.NewBuffer(buf))
	if err == nil {
		t.Error("NewRemainingReader() returned no error, though we wanted one")
	}
}

func TestRemainginReader_Read(t *testing.T) {
	var buf = []byte("hello world!")
	reader, err := NewRemainingReader(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("NewRemainginReader() returned unexpected error: %v", err)
	}

	var buf2 = make([]byte, 5)
	n, err := reader.Read(buf2)
	if err != nil {
		t.Errorf("RemainingReader.Len() returned unexpected error: %v", err)
	}
	if n != len(buf2) {
		t.Errorf("RemainginReader.Len() only read %v bytes, wanted %v bytes", n, buf2)
	}
	if reader.Len() != int64(len(buf)-len(buf2)) {
		t.Errorf(
			"RemainingReader.Len() = %v after reading %v bytes, wanted %v",
			reader.Len(),
			len(buf2),
			len(buf)-len(buf2),
		)
	}

	var buf3 = make([]byte, 10)
	n, err = reader.Read(buf3)
	if err != nil {
		t.Errorf("RemainingReader.Read() returned unexepected error %v", err)
	}
	if n != len(buf)-len(buf2) {
		t.Errorf("RemainingReader.Read() read %v bytes, wanted %v", n, len(buf)-len(buf2))
	}
}
