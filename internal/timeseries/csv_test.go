package timeseries

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/innosat-mats/rac-extract-payload/internal/awstools"
)

func getTestFile() (*CSV, *os.File, error) {
	file, err := ioutil.TempFile("", "innosat-mats-test")
	if err != nil {
		return nil, file, err
	}
	return &CSV{writer: file, csvWriter: csv.NewWriter(file)}, file, err

}

func Test_CSV_Close_WithFile(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	csv.Close()
	buf := make([]byte, 10)
	_, err = file.Read(buf)
	if err == nil {
		t.Error("CSV.Close(), didn't Close file")
	}
}
func Test_CSV_Close_Timeseries(t *testing.T) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("localhost"))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(config)
	upload := manager.NewUploader(s3Client)
	var idxUp = 0
	var uploader = func(uploader *manager.Uploader, key string, bodyBuffer io.Reader) {
		idxUp++
	}
	key := "test"
	ts := awstools.NewTimeseries(uploader, upload, key)
	csv := NewCSV(ts, key)
	if idxUp != 0 {
		t.Errorf("Expected 0 uploads before CSV.Close(), found %v", idxUp)
	}
	csv.Close()
	if idxUp != 1 {
		t.Errorf("Expected 1 upload after CSV.Close(), found %v", idxUp)
	}
}

func Test_CSV_SetSpecifications(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	err = csv.SetSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("CSV.SetSpecifications() = %v, wanted %v", err, nil)
	}
	if !csv.HasSpec {
		t.Errorf("CSV.SetSpecifications() resulted in CSV.HasSpec = %v, wanted %v", csv.HasSpec, true)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.SetSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("CSV.SetSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_CSV_SetSpecifications_no_run_twice(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	// First write should be OK
	err = csv.SetSpecifications([]string{"test", "me"})
	if err != nil {
		t.Errorf("First CSV.SetSpecifications() = %v, wanted %v", err, nil)
	}

	// Second write should be NOK
	err = csv.SetSpecifications([]string{"test", "me"})
	if err == nil {
		t.Errorf("Second CSV.SetSpecifications() = %v, wanted an error", err)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.SetSpecifications() output file could not be located %v", err)
	}
	var want string = "test,me\n"
	if string(content) != want {
		t.Errorf("CSV.SetSpecifications() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_CSV_SetHeaderRow_requires_SetSpecifications(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	err = csv.SetHeaderRow([]string{"Hello", "World"})
	if err == nil {
		t.Errorf("CSV.SetHeaderRow() = %v, wanted an error", err)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("CSV.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_CSV_SetHeaderRow(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	csv.SetSpecifications([]string{"test", "me"})
	err = csv.SetHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("CSV.SetHeaderRow() = %v, wanted %v", err, nil)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("CSV.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}

}

func Test_CSV_SetHeaderRow_only_one_header(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	csv.SetSpecifications([]string{"test", "me"})
	// First Header
	err = csv.SetHeaderRow([]string{"Hello", "World"})
	if err != nil {
		t.Errorf("CSV.SetHeaderRow() = %v, wanted %v", err, nil)
	}
	// Second Header
	err = csv.SetHeaderRow([]string{"World", "World"})
	if err == nil {
		t.Errorf("CSV.SetHeaderRow() = %v, wanted an error", err)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.SetHeaderRow() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\n"
	if string(content) != want {
		t.Errorf("CSV.SetHeaderRow() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_TimeserisCSV_WriteData_requires_spec_and_head(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	err = csv.WriteData([]string{"Test", "1"})
	if err == nil {
		t.Errorf("CSV.WriteData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Specifications and/or") {
		t.Errorf("CSV.WriteData() = %v, wanted error to start with 'Specifications and/or'", err)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.WriteData() output file could not be located %v", err)
	}
	var want string = ""
	if string(content) != want {
		t.Errorf("CSV.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_CSV_WriteData(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	csv.SetSpecifications([]string{"test", "me"})
	csv.SetHeaderRow([]string{"Hello", "World"})
	err = csv.WriteData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("CSV.WriteData() = %v, wanted %v", err, nil)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.WriteData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\n"
	if string(content) != want {
		t.Errorf("CSV.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}

func Test_CSV_WriteData_rejects_bad_columned_row(t *testing.T) {
	csv, file, err := getTestFile()
	defer os.Remove(file.Name())
	if err != nil {
		t.Errorf("CSV fixture could not setup: %v", err)
	}
	csv.SetSpecifications([]string{"test", "me"})
	csv.SetHeaderRow([]string{"Hello", "World"})
	// Good
	err = csv.WriteData([]string{"Test", "1"})
	if err != nil {
		t.Errorf("CSV.WriteData() = %v, wanted %v", err, nil)
	}
	// Bad
	err = csv.WriteData([]string{"Test", "1", "2"})
	if err == nil {
		t.Errorf("CSV.WriteData() = %v, wanted an error", err)
	} else if !strings.HasPrefix(err.Error(), "Irregular column") {
		t.Errorf("CSV.WriteData() = %v, wanted error starting with 'Irregular column'", err)
	}
	// Good again
	err = csv.WriteData([]string{"Test", "2"})
	if err != nil {
		t.Errorf("CSV.WriteData() = %v, wanted %v", err, nil)
	}
	csv.Close()
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("CSV.WriteData() output file could not be located %v", err)
	}
	var want string = "test,me\nHello,World\nTest,1\nTest,2\n"
	if string(content) != want {
		t.Errorf("CSV.WriteData() output file content '%v',' wanted '%v'", string(content), want)
	}
}
