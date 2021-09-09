// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/core/state/api"
	josh "github.com/deref/exo/internal/josh/client"
)

type Store struct {
	client *josh.Client
}

var _ api.Store = (*Store)(nil)

func GetStore(client *josh.Client) *Store {
	return &Store{
		client: client,
	}
}

func (c *Store) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (output *api.DescribeWorkspacesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-workspaces", input, &output)
	return
}

func (c *Store) AddWorkspace(ctx context.Context, input *api.AddWorkspaceInput) (output *api.AddWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "add-workspace", input, &output)
	return
}

func (c *Store) RemoveWorkspace(ctx context.Context, input *api.RemoveWorkspaceInput) (output *api.RemoveWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "remove-workspace", input, &output)
	return
}

func (c *Store) ResolveWorkspace(ctx context.Context, input *api.ResolveWorkspaceInput) (output *api.ResolveWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "resolve-workspace", input, &output)
	return
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
