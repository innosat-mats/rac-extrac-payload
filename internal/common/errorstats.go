package common

import (
	"fmt"
	"sort"
	"strconv"
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

func max(x, y uint) uint {
	if x < y {
		return y
	}
	return x
}

// Summarize ...
func (stats *ErrorStats) Summarize() string {
	errs := make([]struct {
		string
		uint
	}, 0)
	if stats.errorsCount == 0 {
		return fmt.Sprintf(
			"\nStatistics\n\nTotal Errors:\t%v\nTotal Packages:\t%v\n",
			stats.errorsCount,
			stats.total,
		)
	}
	var indent uint = 5 // "Count" is 5 characters
	for key, value := range stats.errors {
		indent = max(indent, uint(len(strconv.Itoa(int(value)))))
		errs = append(errs, struct {
			string
			uint
		}{key, value})
	}
	indent += 3 // Spacing to next column
	sort.SliceStable(errs, func(i, j int) bool {
		return errs[i].uint > errs[j].uint
	})
	lines := make([]string, 0)
	for _, item := range errs {
		lines = append(lines, fmt.Sprintf("%-*v%s", indent, item.uint, item.string))
	}

	return fmt.Sprintf(
		"\nStatistics\n\n%-*vError Message\n%s\n\nTotal Errors:\t%v\nTotal Packages:\t%v\n",
		indent,
		"Count",
		strings.Join(lines, "\n"),
		stats.errorsCount,
		stats.total,
	)
}
