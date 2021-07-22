package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/manifest"
)

func Import(r io.Reader) (*manifest.Manifest, error) {
	procfile, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return Convert(procfile)
}

func Convert(comp *Compose) (*manifest.Manifest, error) {
	// TODO: convert compose to exo.
	return &manifest.Manifest{}, nil
}
