package config

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TelemetryConfig struct {
	Disable bool
}

type Config struct {
	HomeDir string
	BinDir  string
	VarDir  string
	RunDir  string

	Telemetry TelemetryConfig
}

type contextKey string

var configContextKey = contextKey("config")

func LoadDefault(cfg *Config) error {
	if cfg.HomeDir == "" {
		cfg.HomeDir = exoHome()
	}

	configFile := os.Getenv("EXO_CONFIG")
	if configFile == "" {
		configFile = filepath.Join(cfg.HomeDir, "config.toml")
		if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
			return err
		}
	}

	if err := loadFromFile(configFile, cfg); err != nil {
		return err
	}

	if cfg.BinDir == "" {
		cfg.BinDir = filepath.Join(cfg.HomeDir, "bin")
	}
	if cfg.RunDir == "" {
		cfg.RunDir = filepath.Join(cfg.HomeDir, "run")
	}
	if cfg.VarDir == "" {
		cfg.VarDir = filepath.Join(cfg.HomeDir, "var")
	}

	return nil
}

func MustLoadDefault(cfg *Config) {
	if err := LoadDefault(cfg); err != nil {
		panic(err)
	}
}

//go:embed defaultconfig.toml
var defaultConfig []byte

func loadFromFile(filePath string, cfg *Config) error {
	_, err := toml.DecodeFile(filePath, cfg)
	if os.IsNotExist(err) {
		if err = os.WriteFile(filePath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("creating initial config file: %w", err)
		}
	}
	return err
}

func WithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configContextKey, cfg)
}

func GetConfig(ctx context.Context) (*Config, bool) {
	cfg, ok := ctx.Value(configContextKey).(*Config)
	return cfg, ok
}
