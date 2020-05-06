package aez

import (
	"reflect"
	"testing"
)

func Test_stringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Has it", args{"Test", []string{"Not", "a", "Test"}}, true},
		{"Has it not", args{"Test", []string{"Not", "a", "Success"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csvHeader(t *testing.T) {
	myStruct := struct {
		TEST int
		ME   int
		NOT  bool
	}{}
	tests := []struct {
		name    string
		exclude []string
		want    []string
	}{
		{
			"Default use",
			[]string{},
			[]string{"TEST", "ME", "NOT"},
		},
		{
			"Exclusion use",
			[]string{"TEST", "NOT"},
			[]string{"ME"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := csvHeader(myStruct, tt.exclude...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("csvHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
