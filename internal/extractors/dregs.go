package extractors

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

const secondsToNano int64 = 1e9
const MaxDeviationNanos int64 = 30 * secondsToNano               // maximum deviation in ns between dregs and packet
var ErrNoDregsPath error = errors.New("no dregs path specified") // error returned when dregs path is unset

// Dregs The Dregs struct is used for reading and writing incomplete multi
//
//	packet data from .dregs files in order to repair split multi packets
//	between batch runs. Dregs is short for "Data Remaining after Extracting
//	Group of Source packets"
type Dregs struct {
	Path    string // Path to dregs directory
	MaxDiff int64  // Maximum deviation allowed for match [ns]
}

func (dregs *Dregs) getDregsFileName(data common.DataRecord) string {
	return fmt.Sprintf(
		"%v/%v.dregs",
		dregs.Path,
		data.TMHeader.Nanoseconds(),
	)
}

// DumpDregs Write buffer dregs file to specified directory
func (dregs *Dregs) DumpDregs(data common.DataRecord) error {
	if dregs.Path == "" {
		return ErrNoDregsPath
	}

	err := os.MkdirAll(dregs.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed creating %v: %v", dregs.Path, err)
	}

	dregsName := dregs.getDregsFileName(data)
	err = os.WriteFile(dregsName, data.Buffer, 0644)
	if err != nil {
		return fmt.Errorf("failed writing %v: %v", dregsName, err)
	}
	return nil
}

// GetDregs Read buffer dregs and return best match (if any)
func (dregs *Dregs) GetDregs(timestamp int64) ([]byte, error) {
	if dregs.Path == "" {
		return nil, ErrNoDregsPath
	}

	dir := os.DirFS(dregs.Path)
	dregsFiles, err := fs.Glob(dir, "*.dregs")
	if err != nil {
		return nil, fmt.Errorf("failed getting dregs files: %v", err)
	}

	var bestDregs string
	var bestDiff int64 = MaxDeviationNanos
	for _, name := range dregsFiles {
		t, err := strconv.ParseInt(
			strings.TrimSuffix(name, path.Ext(name)),
			10,
			64,
		)
		if err != nil {
			fmt.Printf("could not parse %v as int64: %v\n", name, err)
			continue
		}
		diff := timestamp - t
		if diff > 0 && diff < bestDiff {
			bestDiff = diff
			bestDregs = name
		}
	}

	if bestDregs == "" {
		return nil, fmt.Errorf(
			"found no matching dregs for timestamp %v",
			timestamp,
		)
	}
	return os.ReadFile(fmt.Sprintf("%v/%v", dregs.Path, bestDregs))
}
