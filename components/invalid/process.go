package invalid

import (
	"context"

	core "github.com/deref/exo/core/api"
)

func (p *Provider) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	return nil, p.Err
}

func (p *Provider) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	return nil, p.Err
}

func (p *Provider) Restart(ctx context.Context, input *core.RestartInput) (*core.RestartOutput, error) {
	return nil, p.Err
}
