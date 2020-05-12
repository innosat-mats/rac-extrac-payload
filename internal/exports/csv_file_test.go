package exports

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func getTestFile() (*csvFile, *os.File, error) {
	file, err := ioutil.TempFile("", "innosat-mats-test")
	if err != nil {
		return nil, file, err
	}
	return &csvFile{File: file, Writer: csv.NewWriter(file)}, file, err

}

func Test_csvFile_close(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.close()
	buf := make([]byte, 10)
	_, err = file.Read(buf)
	if err == nil {
		t.Error("csvFile.close(), didn't close file")
	}
}

func Test_csvFile_setSpecifications(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.setSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("csvFile.setSpecifications() = %v, wanted %v", err, nil)
	}
	if !csvFile.HasSpec {
		t.Errorf("csvFile.setSpecifications() resulted in csvFile.HasSpec = %v, wanted %v", csvFile.HasSpec, true)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.setSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("csvFile.setSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_setSpecifications_no_run_twice(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	// First write should be OK
	err = csvFile.setSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("First csvFile.setSpecifications() = %v, wanted %v", err, nil)
	}

	// Second write should be NOK
	err = csvFile.setSpecifications([]string{"test", "me"})
	if err == nil {
		t.Errorf("Second csvFile.setSpecifications() = %v, wanted an error", err)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.setSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("csvFile.setSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_setHeaderRow_requires_setSpecifications(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.setHeaderRow([]string{"Hello", "World"})
	if err == nil {
		t.Errorf("csvFile.setHeaderRow() = %v, wanted an error", err)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.setHeaderRow() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("csvFile.setHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_csvFile_setHeaderRow(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.setSpecifications([]string{"test", "me"})
	err = csvFile.setHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("csvFile.setHeaderRow() = %v, wanted %v", err, nil)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.setHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("csvFile.setHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_csvFile_setHeaderRow_only_one_header(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.setSpecifications([]string{"test", "me"})
	// First Header
	err = csvFile.setHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("csvFile.setHeaderRow() = %v, wanted %v", err, nil)
	}
	// Second Header
	err = csvFile.setHeaderRow([]string{"World", "World"})
	if err == nil {
		t.Errorf("csvFile.setHeaderRow() = %v, wanted an error", err)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.setHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("csvFile.setHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData_requires_spec_and_head(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.writeData([]string{"Test", "1"})
	if err == nil {
		t.Errorf("csvFile.writeData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Specifications and/or") {
		t.Errorf("csvFile.writeData() = %v, wanted error to start with 'Specifications and/or'", err)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.writeData() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("csvFile.writeData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.setSpecifications([]string{"test", "me"})
	csvFile.setHeaderRow([]string{"Hello", "World"})
	err = csvFile.writeData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("csvFile.writeData() = %v, wanted %v", err, nil)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.writeData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\n"
	if string(content) != want {
		t.Errorf("csvFile.writeData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData_rejects_bad_columned_row(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.setSpecifications([]string{"test", "me"})
	csvFile.setHeaderRow([]string{"Hello", "World"})
	// Good
	err = csvFile.writeData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("csvFile.writeData() = %v, wanted %v", err, nil)
	}
	// Bad
	err = csvFile.writeData([]string{"Test", "1", "2"})
	if err == nil {
		t.Errorf("csvFile.writeData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Irregular column") {
		t.Errorf("csvFile.writeData() = %v, wanted error starting with 'Irregular column'", err)
	}
	// Good again
	err = csvFile.writeData([]string{"Test", "2"})
	if err != nil {
		t.Errorf("csvFile.writeData() = %v, wanted %v", err, nil)
	}
	csvFile.close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.writeData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\nTest,2\n"
	if string(content) != want {
		t.Errorf("csvFile.writeData() output file content '%v',' wanted '%v'", string(content), want)
	}
}
