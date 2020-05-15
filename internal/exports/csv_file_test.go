package exports

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func getTestFile() (*CsvFile, *os.File, error) {
	file, err := ioutil.TempFile("", "innosat-mats-test")
	if err != nil {
		return nil, file, err
	}
	return &CsvFile{file: file, writer: csv.NewWriter(file)}, file, err

}

func Test_csvFile_close(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.Close()
	buf := make([]byte, 10)
	_, err = file.Read(buf)
	if err == nil {
		t.Error("csvFile.Close(), didn't Close file")
	}
}

func Test_csvFile_setSpecifications(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.SetSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("csvFile.SetSpecifications() = %v, wanted %v", err, nil)
	}
	if !csvFile.HasSpec {
		t.Errorf("csvFile.SetSpecifications() resulted in csvFile.HasSpec = %v, wanted %v", csvFile.HasSpec, true)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.SetSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("csvFile.SetSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_setSpecifications_no_run_twice(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	// First write should be OK
	err = csvFile.SetSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("First csvFile.SetSpecifications() = %v, wanted %v", err, nil)
	}

	// Second write should be NOK
	err = csvFile.SetSpecifications([]string{"test", "me"})
	if err == nil {
		t.Errorf("Second csvFile.SetSpecifications() = %v, wanted an error", err)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.SetSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("csvFile.SetSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_setHeaderRow_requires_setSpecifications(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.SetHeaderRow([]string{"Hello", "World"})
	if err == nil {
		t.Errorf("csvFile.SetHeaderRow() = %v, wanted an error", err)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("csvFile.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_csvFile_setHeaderRow(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.SetSpecifications([]string{"test", "me"})
	err = csvFile.SetHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("csvFile.SetHeaderRow() = %v, wanted %v", err, nil)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("csvFile.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_csvFile_setHeaderRow_only_one_header(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.SetSpecifications([]string{"test", "me"})
	// First Header
	err = csvFile.SetHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("csvFile.SetHeaderRow() = %v, wanted %v", err, nil)
	}
	// Second Header
	err = csvFile.SetHeaderRow([]string{"World", "World"})
	if err == nil {
		t.Errorf("csvFile.SetHeaderRow() = %v, wanted an error", err)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("csvFile.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData_requires_spec_and_head(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	err = csvFile.WriteData([]string{"Test", "1"})
	if err == nil {
		t.Errorf("csvFile.WriteData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Specifications and/or") {
		t.Errorf("csvFile.WriteData() = %v, wanted error to start with 'Specifications and/or'", err)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.WriteData() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("csvFile.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.SetSpecifications([]string{"test", "me"})
	csvFile.SetHeaderRow([]string{"Hello", "World"})
	err = csvFile.WriteData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("csvFile.WriteData() = %v, wanted %v", err, nil)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.WriteData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\n"
	if string(content) != want {
		t.Errorf("csvFile.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_csvFile_writeData_rejects_bad_columned_row(t *testing.T) {
	csvFile, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("csvFile fixture could not setup: %v", err)
	}
	csvFile.SetSpecifications([]string{"test", "me"})
	csvFile.SetHeaderRow([]string{"Hello", "World"})
	// Good
	err = csvFile.WriteData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("csvFile.WriteData() = %v, wanted %v", err, nil)
	}
	// Bad
	err = csvFile.WriteData([]string{"Test", "1", "2"})
	if err == nil {
		t.Errorf("csvFile.WriteData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Irregular column") {
		t.Errorf("csvFile.WriteData() = %v, wanted error starting with 'Irregular column'", err)
	}
	// Good again
	err = csvFile.WriteData([]string{"Test", "2"})
	if err != nil {
		t.Errorf("csvFile.WriteData() = %v, wanted %v", err, nil)
	}
	csvFile.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("csvFile.WriteData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\nTest,2\n"
	if string(content) != want {
		t.Errorf("csvFile.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}
