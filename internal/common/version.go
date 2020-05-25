package common

import "fmt"

// Version is the version of the source code
var Version string

// Head is the short commit id of head
var Head string

// Buildtime is the time of the build
var Buildtime string

// FullVersion returns a complete version string
func FullVersion() string {
	if Version == "" && Head == "" && Buildtime == "" {
		return "Test Build"
	}

	return fmt.Sprintf("%v (%v) @ %v", Version, Head, Buildtime)
}
