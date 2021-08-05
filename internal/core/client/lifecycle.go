// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	josh "github.com/deref/exo/internal/josh/client"
)

type Lifecycle struct {
	client *josh.Client
}

var _ api.Lifecycle = (*Lifecycle)(nil)

func GetLifecycle(client *josh.Client) *Lifecycle {
	return &Lifecycle{
		client: client,
	}
}

func (c *Lifecycle) Initialize(ctx context.Context, input *api.InitializeInput) (output *api.InitializeOutput, err error) {
	err = c.client.Invoke(ctx, "initialize", input, &output)
	return
}

func (c *Lifecycle) Update(ctx context.Context, input *api.UpdateInput) (output *api.UpdateOutput, err error) {
	err = c.client.Invoke(ctx, "update", input, &output)
	return
}

func (c *Lifecycle) Refresh(ctx context.Context, input *api.RefreshInput) (output *api.RefreshOutput, err error) {
	err = c.client.Invoke(ctx, "refresh", input, &output)
	return
}

func (c *Lifecycle) Dispose(ctx context.Context, input *api.DisposeInput) (output *api.DisposeOutput, err error) {
	err = c.client.Invoke(ctx, "dispose", input, &output)
	return
}
