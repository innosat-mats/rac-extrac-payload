package exports

import (
	"bytes"
	"reflect"
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
		want      string
	}{
		{"Prints nothing when not asked to", args{false}, innerArgs{common.DataRecord{}}, ""},
		{
			"Prints something when asked to",
			args{true},
			innerArgs{common.DataRecord{}},
			"{Origin:{Name: ProcessingDate:0001-01-01 00:00:00 +0000 UTC} RamsesHeader:{Synch:0 Length:0 Port:0 Type:0 Secure:0 Time:0 Date:0} RamsesSecure:{IPAddress:0 Port:0 Seq:0 Retransmission:0 Ack:0 _:0} SourceHeader:{PacketID:0 PacketSequenceControl:0 PacketLength:0} TMHeader:{PUS:0 ServiceType:0 ServiceSubType:0 CUCTimeSeconds:0 CUCTimeFraction:0} SID: Data:<nil> Error:<nil> Buffer:[]}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			callback, _ := StdoutCallbackFactory(buf, tt.args.writeTimeseries)
			callback(tt.innerArgs.dataRecord)
			if got := buf.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StdoutCallbackFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}
