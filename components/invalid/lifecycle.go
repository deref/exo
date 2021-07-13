package invalid

import (
	"context"

	"github.com/deref/exo/kernel/api"
)

type Lifecycle struct {
	Err error
}

func (lc *Lifecycle) Initialize(ctx context.Context, input *api.InitializeInput) (*api.InitializeOutput, error) {
	return nil, lc.Err
}

func (lc *Lifecycle) Update(context.Context, *api.UpdateInput) (*api.UpdateOutput, error) {
	return nil, lc.Err
}

func (lc *Lifecycle) Refresh(context.Context, *api.RefreshInput) (*api.RefreshOutput, error) {
	return nil, lc.Err
}

func (lc *Lifecycle) Dispose(context.Context, *api.DisposeInput) (*api.DisposeOutput, error) {
	return nil, lc.Err
}
