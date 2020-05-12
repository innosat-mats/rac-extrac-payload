package exports

import (
	"bytes"
	"strings"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_StdoutCallbackFactory(t *testing.T) {
	type args struct {
		writeTimeseries bool
	}
	type innerArgs struct {
		dataRecord common.DataRecord
	}
	tests := []struct {
		name      string
		args      args
		innerArgs innerArgs
		wantLines int
	}{
		{"Prints nothing when not asked to", args{false}, innerArgs{common.DataRecord{}}, 0},
		{
			"Prints something when asked to",
			args{true},
			innerArgs{common.DataRecord{}},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			callback, _ := StdoutCallbackFactory(buf, tt.args.writeTimeseries)
			callback(tt.innerArgs.dataRecord)
			if got := strings.Count(buf.String(), "\n"); got != tt.wantLines {
				t.Errorf("StdoutCallbackFactory() is %v lines, but want %v", got, tt.wantLines)
			}
		})
	}
}
