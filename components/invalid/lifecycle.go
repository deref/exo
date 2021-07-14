package invalid

import (
	"context"

	"github.com/deref/exo/core"
)

func (p *Provider) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	return nil, p.Err
}

func (p *Provider) Update(context.Context, *core.UpdateInput) (*core.UpdateOutput, error) {
	return nil, p.Err
}

func (p *Provider) Refresh(context.Context, *core.RefreshInput) (*core.RefreshOutput, error) {
	return nil, p.Err
}

func (p *Provider) Dispose(context.Context, *core.DisposeInput) (*core.DisposeOutput, error) {
	return nil, p.Err
}
