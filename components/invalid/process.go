package invalid

import (
	"context"

	"github.com/deref/exo/core"
)

func (p *Provider) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	return nil, p.Err
}

func (p *Provider) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	return nil, p.Err
}
