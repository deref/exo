// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/api"
)

type Kernel struct {
	client *josh.Client
}

var _ api.Kernel = (*Kernel)(nil)

func NewKernel(client *josh.Client) *Kernel {
	return &Kernel{
		client: client,
	}
}

func (c *Kernel) CreateWorkspace(ctx context.Context, input *api.CreateWorkspaceInput) (output *api.CreateWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "create-workspace", input, &output)
	return
}

func (c *Kernel) ForgetWorkspace(ctx context.Context, input *api.ForgetWorkspaceInput) (output *api.ForgetWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "forget-workspace", input, &output)
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
