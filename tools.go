// +build tools

// Package tools keeps track of our development tools.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "golang.org/x/lint" // Linting
)
