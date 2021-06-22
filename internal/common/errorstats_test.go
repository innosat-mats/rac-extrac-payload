package common

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestErrorStats_Register_and_Summarize(t *testing.T) {
	tests := []struct {
		name            string
		errs            []error
		wantTotal       uint
		wantTotalErrors uint
		wantErrors      []string
	}{
		{
			"Empty",
			[]error{},
			0,
			0,
			[]string{},
		},
		{
			"Some Errors",
			[]error{
				errors.New("test"),
				nil,
				io.ErrUnexpectedEOF,
				nil,
				io.ErrUnexpectedEOF,
			},
			5,
			3,
			[]string{
				"2\tunexpected EOF",
				"1\ttest",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewErrorStats()
			for _, err := range tt.errs {
				stats.Register(err)
			}
			want := fmt.Sprintf(
				"\nStatistics\n\nCount\tError Message\n%s\n\nTotal Errors:\t%v\nTotal Packages:\t%v\n",
				strings.Join(tt.wantErrors, "\n"),
				tt.wantTotalErrors,
				tt.wantTotal,
			)
			if got := stats.Summarize(); got != want {
				t.Errorf("ErrorStats.Summarize() = %v, want %v", got, want)
			}
		})
	}
}
