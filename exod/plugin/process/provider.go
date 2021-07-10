package process

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/deref/exo/exod/api"
	"github.com/mitchellh/mapstructure"
)

type Provider struct{}

type spec struct {
	Command   string
	Arguments []string
}

func (provider *Provider) Create(ctx context.Context, input *api.CreateInput) (*api.CreateOutput, error) {
	var spec spec
	if err := mapstructure.Decode(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("decoding mapstructure: %w", err)
	}

	cmd := exec.Command(spec.Command, spec.Arguments...)
	// TODO: handle stdio.
	// TODO: clear/reset environment.
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting: %w", err)
	}

	var output api.CreateOutput
	output.IRI = fmt.Sprintf("process:localhost:%d", cmd.Process.Pid)
	return &output, nil
}
