package invalid

import (
	"context"

	"github.com/deref/exo/kernel/api"
)

func (p *Provider) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	return nil, p.Err
}

func (p *Provider) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	return nil, p.Err
}
