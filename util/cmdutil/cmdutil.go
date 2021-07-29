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
	ConfigFile   string // Contains global configuration
}

func MustMakeDirectories() *KnownPaths {
	var paths KnownPaths
	mkdir := func(out *string, path string) {
		if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
			Fatalf("making %q: %w", path, err)
		}
		*out = path
	}

	mkdir(&paths.ExoDir, config.ExoHome())
	mkdir(&paths.BinDir, config.ExoBin())
	mkdir(&paths.VarDir, config.ExoVar())
	mkdir(&paths.RunDir, config.ExoRun())

	paths.RunStateFile = filepath.Join(paths.RunDir, "exod.json")
	paths.ConfigFile = config.ConfigFile()

	return &paths
}

func GetAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	return "localhost:" + port
}

func MustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
