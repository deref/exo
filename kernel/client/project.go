// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/api"
)

type Project struct {
	client *josh.Client
}

var _ api.Project = (*Project)(nil)

func NewProject(client *josh.Client) *Project {
	return &Project{
		client: client,
	}
}

func (c *Project) Delete(ctx context.Context, input *api.DeleteInput) (output *api.DeleteOutput, err error) {
	err = c.client.Invoke(ctx, "delete", input, &output)
	return
}

func (c *Project) Apply(ctx context.Context, input *api.ApplyInput) (output *api.ApplyOutput, err error) {
	err = c.client.Invoke(ctx, "apply", input, &output)
	return
}

func (c *Project) ApplyProcfile(ctx context.Context, input *api.ApplyProcfileInput) (output *api.ApplyProcfileOutput, err error) {
	err = c.client.Invoke(ctx, "apply-procfile", input, &output)
	return
}

func (c *Project) Refresh(ctx context.Context, input *api.RefreshInput) (output *api.RefreshOutput, err error) {
	err = c.client.Invoke(ctx, "refresh", input, &output)
	return
}

func (c *Project) Resolve(ctx context.Context, input *api.ResolveInput) (output *api.ResolveOutput, err error) {
	err = c.client.Invoke(ctx, "resolve", input, &output)
	return
}

func (c *Project) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (output *api.DescribeComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-components", input, &output)
	return
}

func (c *Project) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (output *api.CreateComponentOutput, err error) {
	err = c.client.Invoke(ctx, "create-component", input, &output)
	return
}

func (c *Project) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (output *api.UpdateComponentOutput, err error) {
	err = c.client.Invoke(ctx, "update-component", input, &output)
	return
}

func (c *Project) RefreshComponent(ctx context.Context, input *api.RefreshComponentInput) (output *api.RefreshComponentOutput, err error) {
	err = c.client.Invoke(ctx, "refresh-component", input, &output)
	return
}

func (c *Project) DisposeComponent(ctx context.Context, input *api.DisposeComponentInput) (output *api.DisposeComponentOutput, err error) {
	err = c.client.Invoke(ctx, "dispose-component", input, &output)
	return
}

func (c *Project) DeleteComponent(ctx context.Context, input *api.DeleteComponentInput) (output *api.DeleteComponentOutput, err error) {
	err = c.client.Invoke(ctx, "delete-component", input, &output)
	return
}

func (c *Project) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (output *api.DescribeLogsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-logs", input, &output)
	return
}

func (c *Project) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = c.client.Invoke(ctx, "get-events", input, &output)
	return
}

func (c *Project) Start(ctx context.Context, input *api.StartInput) (output *api.StartOutput, err error) {
	err = c.client.Invoke(ctx, "start", input, &output)
	return
}

func (c *Project) Stop(ctx context.Context, input *api.StopInput) (output *api.StopOutput, err error) {
	err = c.client.Invoke(ctx, "stop", input, &output)
	return
}

func (c *Project) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (output *api.DescribeProcessesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-processes", input, &output)
	return
}
