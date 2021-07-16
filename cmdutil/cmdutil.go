package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func Fatalf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf(format, v...))
	os.Exit(1)
}

func MustVarDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("getting home directory: %w", err))
	}
	varDir := filepath.Join(homeDir, ".exo", "var")
	if err := os.MkdirAll(varDir, 0700); err != nil {
		panic(fmt.Errorf("mk var dir: %w", err))
	}
	return varDir
}
