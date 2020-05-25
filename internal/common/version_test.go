package common

import "testing"

func TestFullVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		head      string
		buildtime string
		want      string
	}{
		{name: "Unset causes Test Build", want: "Test Build"},
		{
			"Returns full version",
			"v1.0.0",
			"12345",
			"2020-02-02T22:22:22Z",
			"v1.0.0 (12345) @ 2020-02-02T22:22:22Z",
		},
	}
	var version = Version
	var head = Head
	var buildtime = Buildtime

	cleanup := func() {
		Version = version
		Head = head
		Buildtime = buildtime
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Version = tt.version
			Head = tt.head
			Buildtime = tt.buildtime
			defer cleanup()
			if got := FullVersion(); got != tt.want {
				t.Errorf("FullVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
