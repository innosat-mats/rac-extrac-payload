package common

import (
	"fmt"
	"sort"
	"strings"
)

// ErrorStats ...
type ErrorStats struct {
	total       uint // errors and not errors
	errorsCount uint
	errors      map[string]uint
}

// NewErrorStats ...
func NewErrorStats() ErrorStats {
	return ErrorStats{errors: make(map[string]uint)}
}

// Register a new error or not error occurnace
func (stats *ErrorStats) Register(err error) {
	stats.total++
	if err != nil {
		stats.errorsCount++
		stats.errors[err.Error()]++
	}

}

// Summarize ...
func (stats *ErrorStats) Summarize() string {
	errs := make([]struct {
		string
		uint
	}, 0)
	for key, value := range stats.errors {
		errs = append(errs, struct {
			string
			uint
		}{key, value})
	}
	sort.SliceStable(errs, func(i, j int) bool {
		return errs[i].uint > errs[j].uint
	})
	lines := make([]string, 0)
	for _, item := range errs {
		lines = append(lines, fmt.Sprintf("%v\t%s", item.uint, item.string))
	}

	return fmt.Sprintf(
		"\nStatistics\n\nCount\tError Message\n%s\n\nTotal Errors:\t%v\nTotal Packages:\t%v\n",
		strings.Join(lines, "\n"),
		stats.errorsCount,
		stats.total,
	)
}
