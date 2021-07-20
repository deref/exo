// +build bundle

package main

import "github.com/deref/exo/components/process"

var fifofumConfig = process.FifofumConfig{
	Path: "exo",
	Args: []string{"fifofum"},
}
