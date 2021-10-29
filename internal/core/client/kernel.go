// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	josh "github.com/deref/exo/internal/josh/client"
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

func (c *Kernel) AuthEsv(ctx context.Context, input *api.AuthEsvInput) (output *api.AuthEsvOutput, err error) {
	err = c.client.Invoke(ctx, "auth-esv", input, &output)
	return
}

func (c *Kernel) UnauthEsv(ctx context.Context, input *api.UnauthEsvInput) (output *api.UnauthEsvOutput, err error) {
	err = c.client.Invoke(ctx, "unauth-esv", input, &output)
	return
}

func (c *Kernel) GetEsvUser(ctx context.Context, input *api.GetEsvUserInput) (output *api.GetEsvUserOutput, err error) {
	err = c.client.Invoke(ctx, "get-esv-user", input, &output)
	return
}

func (c *Kernel) CreateProject(ctx context.Context, input *api.CreateProjectInput) (output *api.CreateProjectOutput, err error) {
	err = c.client.Invoke(ctx, "create-project", input, &output)
	return
}

func (c *Kernel) DescribeTemplates(ctx context.Context, input *api.DescribeTemplatesInput) (output *api.DescribeTemplatesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-templates", input, &output)
	return
}

func (c *Kernel) CreateWorkspace(ctx context.Context, input *api.CreateWorkspaceInput) (output *api.CreateWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "create-workspace", input, &output)
	return
}

func (c *Kernel) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (output *api.DescribeWorkspacesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-workspaces", input, &output)
	return
}

func (c *Kernel) ResolveWorkspace(ctx context.Context, input *api.ResolveWorkspaceInput) (output *api.ResolveWorkspaceOutput, err error) {
	err = c.client.Invoke(ctx, "resolve-workspace", input, &output)
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

func (c *Kernel) Upgrade(ctx context.Context, input *api.UpgradeInput) (output *api.UpgradeOutput, err error) {
	err = c.client.Invoke(ctx, "upgrade", input, &output)
	return
}

func (c *Kernel) Ping(ctx context.Context, input *api.PingInput) (output *api.PingOutput, err error) {
	err = c.client.Invoke(ctx, "ping", input, &output)
	return
}

func (c *Kernel) Exit(ctx context.Context, input *api.ExitInput) (output *api.ExitOutput, err error) {
	err = c.client.Invoke(ctx, "exit", input, &output)
	return
}

func (c *Kernel) DescribeTasks(ctx context.Context, input *api.DescribeTasksInput) (output *api.DescribeTasksOutput, err error) {
	err = c.client.Invoke(ctx, "describe-tasks", input, &output)
	return
}

func (c *Kernel) GetUserHomeDir(ctx context.Context, input *api.GetUserHomeDirInput) (output *api.GetUserHomeDirOutput, err error) {
	err = c.client.Invoke(ctx, "get-user-home-dir", input, &output)
	return
}

func (c *Kernel) ReadDir(ctx context.Context, input *api.ReadDirInput) (output *api.ReadDirOutput, err error) {
	err = c.client.Invoke(ctx, "read-dir", input, &output)
	return
}
