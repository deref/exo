package main

import (
	"encoding/json"
	"os"

	"github.com/deref/exo/internal/supervise"
	"github.com/deref/exo/internal/util/cmdutil"
)

func main() {
	cfg := &supervise.Config{}
	if err := json.NewDecoder(os.Stdin).Decode(cfg); err != nil {
		cmdutil.Fatalf("reading config from stdin: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		cmdutil.Fatalf("validating config: %v", err)
	}

	supervise.Main(cfg)
}
