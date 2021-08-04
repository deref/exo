package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ClientConfig struct {
	URL string `toml:"url"`
}

type GUIConfig struct {
	Port uint
}

type LogConfig struct {
	SyslogPort uint
}

type TelemetryConfig struct {
	Disable bool
}

type Config struct {
	HomeDir string
	BinDir  string
	VarDir  string
	RunDir  string

	HTTPPort uint `toml:"httpPort"`

	Client    ClientConfig
	GUI       GUIConfig `toml:"gui"`
	Log       LogConfig
	Telemetry TelemetryConfig
}

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

	setDefaults(cfg)

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

func setDefaults(cfg *Config) {
	// Paths
	if cfg.BinDir == "" {
		cfg.BinDir = filepath.Join(cfg.HomeDir, "bin")
	}
	if cfg.RunDir == "" {
		cfg.RunDir = filepath.Join(cfg.HomeDir, "run")
	}
	if cfg.VarDir == "" {
		cfg.VarDir = filepath.Join(cfg.HomeDir, "var")
	}

	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 4000
	}

	// Log
	if cfg.Log.SyslogPort == 0 {
		cfg.Log.SyslogPort = 4500
	}

	// GUI
	if cfg.GUI.Port == 0 {
		cfg.GUI.Port = 3000
	}
}
