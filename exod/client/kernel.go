// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/exod/api"
	josh "github.com/deref/exo/josh/client"
)

type Kernel struct {
	client *josh.Client
}

var _ api.Kernel = (*Kernel)(nil)

func GetKernel(client *josh.Client) *Kernel {
	return &Kernel{
		client: client,
	}
}

func (c *Kernel) CreateWorkspace(ctx context.Context, input *api.CreateWorkspaceInput) (output *api.CreateWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "create-workspace", input, &output)
	return
}

func (c *Kernel) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (output *api.DescribeWorkspacesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-workspaces", input, &output)
	return
}

func (c *Kernel) FindWorkspace(ctx context.Context, input *api.FindWorkspaceInput) (output *api.FindWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "find-workspace", input, &output)
	return
}

func (c *Kernel) Panic(ctx context.Context, input *api.PanicInput) (output *api.PanicOutput, err error) {
	err = c.client.Invoke(ctx, "panic", input, &output)
	return
}

func (c *Kernel) GetVersion(ctx context.Context, input *api.GetVersionInput) (output *api.GetVersionOutput, err error) {
	err = c.client.Invoke(ctx, "get-version", input, &output)
	return
}
