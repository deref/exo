package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/config"
)

func Import(r io.Reader) (*config.Config, error) {
	procfile, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return Convert(procfile)
}

func Convert(comp *Compose) (*config.Config, error) {
	// TODO: convert compose to exo.
	return &config.Config{}, nil
}
