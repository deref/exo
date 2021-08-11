package supervise

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Config struct {
	ComponentID      string
	WorkingDirectory string
	Environment      map[string]string
	SyslogPort       uint
	Program          string
	Arguments        []string
}

func (cfg *Config) Validate() error {
	var errorMessages []string
	if cfg.ComponentID == "" {
		errorMessages = append(errorMessages, "missing ComponentID")
	}
	if cfg.WorkingDirectory == "" {
		errorMessages = append(errorMessages, "missing WorkingDirectory")
	}
	if cfg.SyslogPort == 0 {
		errorMessages = append(errorMessages, "missing SyslogPort")
	}
	if cfg.Program == "" {
		errorMessages = append(errorMessages, "missing Program")
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("invalid supervisor config: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

func MustEncodeConfig(cfg *Config) []byte {
	out, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return out
}
