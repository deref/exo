package process

import (
	"context"

	"github.com/deref/exo/kernel/api"
)

func (p *Provider) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	panic("TODO: process start")
}

func (p *Provider) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	panic("TODO: process stop")
}
