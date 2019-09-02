// +build tools

// A workaround to get the tool packages in the go.mod. See: https://github.com/golang/go/issues/25922
// somehow it doesn't work...
package main

import (
	_ "github.com/mattn/goveralls"
	_ "github.com/rakyll/gotest"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/cover"
)
