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
