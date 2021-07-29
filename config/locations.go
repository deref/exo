package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func exoHome() string {
	if homeFromEnv := os.Getenv("EXO_HOME"); homeFromEnv != "" {
		return homeFromEnv
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("getting home directory: %w", err))
	}

	return filepath.Join(homeDir, ".exo")
}

func getExoPath(homedir, subdir, envOverride string) string {
	if dirFromEnv := os.Getenv(envOverride); dirFromEnv != "" {
		return dirFromEnv
	}
	return filepath.Join(homedir, subdir)
}
