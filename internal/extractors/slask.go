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
const MaxDeviationNanos int64 = 30 * secondsToNano               // maximum deviation in ns between slask and packet
var ErrNoSlaskPath error = errors.New("no slask path specified") // error returned when slask path is unset

// Slask The Slask struct is used for reading and writing incomplete multi
// packet data to help repairing split multi packets between batch runs
type Slask struct {
	Path    string // Path to slask directory
	MaxDiff int64  // Maximum deviation allowed for match [ns]
}

func (slask *Slask) getSlaskFileName(data common.DataRecord) string {
	return fmt.Sprintf(
		"%v/%v.slask",
		slask.Path,
		data.TMHeader.Nanoseconds(),
	)
}

// DumpSlask Write buffer slask file to specified directory
func (slask *Slask) DumpSlask(data common.DataRecord) error {
	if slask.Path == "" {
		return ErrNoSlaskPath
	}

	err := os.MkdirAll(slask.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed creating %v: %v", slask.Path, err)
	}

	slaskName := slask.getSlaskFileName(data)
	err = os.WriteFile(slaskName, data.Buffer, 0644)
	if err != nil {
		return fmt.Errorf("failed writing %v: %v", slaskName, err)
	}
	return nil
}

// GetSlask Read buffer slask and return best match (if any)
func (slask *Slask) GetSlask(timestamp int64) ([]byte, error) {
	if slask.Path == "" {
		return nil, ErrNoSlaskPath
	}

	dir := os.DirFS(slask.Path)
	slaskFiles, err := fs.Glob(dir, "*.slask")
	if err != nil {
		return nil, fmt.Errorf("failed getting slask files: %v", err)
	}

	var bestSlask string
	var bestDiff int64 = MaxDeviationNanos
	for _, name := range slaskFiles {
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
			bestSlask = name
		}
	}

	if bestSlask == "" {
		return nil, fmt.Errorf(
			"found no matching slask for timestamp %v",
			timestamp,
		)
	}
	return os.ReadFile(fmt.Sprintf("%v/%v", slask.Path, bestSlask))
}
