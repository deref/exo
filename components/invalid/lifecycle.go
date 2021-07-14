package invalid

import (
	"context"

	"github.com/deref/exo/kernel/api"
)

func (p *Provider) Initialize(ctx context.Context, input *api.InitializeInput) (*api.InitializeOutput, error) {
	return nil, p.Err
}

func (p *Provider) Update(context.Context, *api.UpdateInput) (*api.UpdateOutput, error) {
	return nil, p.Err
}

func (p *Provider) Refresh(context.Context, *api.RefreshInput) (*api.RefreshOutput, error) {
	return nil, p.Err
}

func (p *Provider) Dispose(context.Context, *api.DisposeInput) (*api.DisposeOutput, error) {
	return nil, p.Err
}
