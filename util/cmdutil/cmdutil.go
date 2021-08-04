package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deref/exo/config"
)

func Warnf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	Warnf(format, v...)
	os.Exit(1)
}

func Fatal(err error) {
	Fatalf("%v", err)
}

type KnownPaths struct {
	ExoDir string // Exo home directory.
	BinDir string // Binaries.
	VarDir string // Durable state.
	RunDir string // Volatile state.

	RunStateFile string // Contains information about the exo daemon.
}

func MustMakeDirectories(cfg *config.Config) *KnownPaths {
	var paths KnownPaths
	mkdir := func(out *string, path string) {
		if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
			Fatalf("making %q: %w", path, err)
		}
		*out = path
	}

	mkdir(&paths.ExoDir, cfg.HomeDir)
	mkdir(&paths.BinDir, cfg.BinDir)
	mkdir(&paths.VarDir, cfg.VarDir)
	mkdir(&paths.RunDir, cfg.RunDir)

	paths.RunStateFile = filepath.Join(paths.RunDir, "exod.json")

	return &paths
}

func GetAddr(cfg *config.Config) string {
	return fmt.Sprintf("localhost:%d", cfg.HTTPPort)
}

func MustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
