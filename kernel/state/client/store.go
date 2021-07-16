// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/state/api"
)

type Store struct {
	client *josh.Client
}

var _ api.Store = (*Store)(nil)

func NewStore(client *josh.Client) *Store {
	return &Store{
		client: client,
	}
}

func (c *Store) Resolve(ctx context.Context, input *api.ResolveInput) (output *api.ResolveOutput, err error) {
	err = c.client.Invoke(ctx, "resolve", input, &output)
	return
}

func (c *Store) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (output *api.DescribeComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-components", input, &output)
	return
}

func (c *Store) AddComponent(ctx context.Context, input *api.AddComponentInput) (output *api.AddComponentOutput, err error) {
	err = c.client.Invoke(ctx, "add-component", input, &output)
	return
}

func (c *Store) PatchComponent(ctx context.Context, input *api.PatchComponentInput) (output *api.PatchComponentOutput, err error) {
	err = c.client.Invoke(ctx, "patch-component", input, &output)
	return
}

func (c *Store) RemoveComponent(ctx context.Context, input *api.RemoveComponentInput) (output *api.RemoveComponentOutput, err error) {
	err = c.client.Invoke(ctx, "remove-component", input, &output)
	return
}
