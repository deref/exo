// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	josh "github.com/deref/exo/internal/josh/client"
)

type Process struct {
	client *josh.Client
}

var _ api.Process = (*Process)(nil)

func GetProcess(client *josh.Client) *Process {
	return &Process{
		client: client,
	}
}

func (c *Process) Start(ctx context.Context, input *api.StartInput) (output *api.StartOutput, err error) {
	err = c.client.Invoke(ctx, "start", input, &output)
	return
}

func (c *Process) Stop(ctx context.Context, input *api.StopInput) (output *api.StopOutput, err error) {
	err = c.client.Invoke(ctx, "stop", input, &output)
	return
}

func (c *Process) Signal(ctx context.Context, input *api.SignalInput) (output *api.SignalOutput, err error) {
	err = c.client.Invoke(ctx, "signal", input, &output)
	return
}

func (c *Process) Restart(ctx context.Context, input *api.RestartInput) (output *api.RestartOutput, err error) {
	err = c.client.Invoke(ctx, "restart", input, &output)
	return
}

type Builder struct {
	client *josh.Client
}

var _ api.Builder = (*Builder)(nil)

func GetBuilder(client *josh.Client) *Builder {
	return &Builder{
		client: client,
	}
}

func (c *Builder) Build(ctx context.Context, input *api.BuildInput) (output *api.BuildOutput, err error) {
	err = c.client.Invoke(ctx, "build", input, &output)
	return
}

type Workspace struct {
	client *josh.Client
}

var _ api.Workspace = (*Workspace)(nil)

func GetWorkspace(client *josh.Client) *Workspace {
	return &Workspace{
		client: client,
	}
}

func (c *Workspace) Start(ctx context.Context, input *api.StartInput) (output *api.StartOutput, err error) {
	err = c.client.Invoke(ctx, "start", input, &output)
	return
}

func (c *Workspace) Stop(ctx context.Context, input *api.StopInput) (output *api.StopOutput, err error) {
	err = c.client.Invoke(ctx, "stop", input, &output)
	return
}

func (c *Workspace) Signal(ctx context.Context, input *api.SignalInput) (output *api.SignalOutput, err error) {
	err = c.client.Invoke(ctx, "signal", input, &output)
	return
}

func (c *Workspace) Restart(ctx context.Context, input *api.RestartInput) (output *api.RestartOutput, err error) {
	err = c.client.Invoke(ctx, "restart", input, &output)
	return
}

func (c *Workspace) Build(ctx context.Context, input *api.BuildInput) (output *api.BuildOutput, err error) {
	err = c.client.Invoke(ctx, "build", input, &output)
	return
}

func (c *Workspace) Describe(ctx context.Context, input *api.DescribeInput) (output *api.DescribeOutput, err error) {
	err = c.client.Invoke(ctx, "describe", input, &output)
	return
}

func (c *Workspace) Destroy(ctx context.Context, input *api.DestroyInput) (output *api.DestroyOutput, err error) {
	err = c.client.Invoke(ctx, "destroy", input, &output)
	return
}

func (c *Workspace) Apply(ctx context.Context, input *api.ApplyInput) (output *api.ApplyOutput, err error) {
	err = c.client.Invoke(ctx, "apply", input, &output)
	return
}

func (c *Workspace) Resolve(ctx context.Context, input *api.ResolveInput) (output *api.ResolveOutput, err error) {
	err = c.client.Invoke(ctx, "resolve", input, &output)
	return
}

func (c *Workspace) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (output *api.DescribeComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-components", input, &output)
	return
}

func (c *Workspace) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (output *api.CreateComponentOutput, err error) {
	err = c.client.Invoke(ctx, "create-component", input, &output)
	return
}

func (c *Workspace) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (output *api.UpdateComponentOutput, err error) {
	err = c.client.Invoke(ctx, "update-component", input, &output)
	return
}

func (c *Workspace) RefreshComponents(ctx context.Context, input *api.RefreshComponentsInput) (output *api.RefreshComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "refresh-components", input, &output)
	return
}

func (c *Workspace) DisposeComponents(ctx context.Context, input *api.DisposeComponentsInput) (output *api.DisposeComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "dispose-components", input, &output)
	return
}

func (c *Workspace) DeleteComponents(ctx context.Context, input *api.DeleteComponentsInput) (output *api.DeleteComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "delete-components", input, &output)
	return
}

func (c *Workspace) GetComponentState(ctx context.Context, input *api.GetComponentStateInput) (output *api.GetComponentStateOutput, err error) {
	err = c.client.Invoke(ctx, "get-component-state", input, &output)
	return
}

func (c *Workspace) SetComponentState(ctx context.Context, input *api.SetComponentStateInput) (output *api.SetComponentStateOutput, err error) {
	err = c.client.Invoke(ctx, "set-component-state", input, &output)
	return
}

func (c *Workspace) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (output *api.DescribeLogsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-logs", input, &output)
	return
}

func (c *Workspace) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = c.client.Invoke(ctx, "get-events", input, &output)
	return
}

func (c *Workspace) StartComponents(ctx context.Context, input *api.StartComponentsInput) (output *api.StartComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "start-components", input, &output)
	return
}

func (c *Workspace) StopComponents(ctx context.Context, input *api.StopComponentsInput) (output *api.StopComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "stop-components", input, &output)
	return
}

func (c *Workspace) SignalComponents(ctx context.Context, input *api.SignalComponentsInput) (output *api.SignalComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "signal-components", input, &output)
	return
}

func (c *Workspace) RestartComponents(ctx context.Context, input *api.RestartComponentsInput) (output *api.RestartComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "restart-components", input, &output)
	return
}

func (c *Workspace) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (output *api.DescribeProcessesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-processes", input, &output)
	return
}

func (c *Workspace) DescribeVolumes(ctx context.Context, input *api.DescribeVolumesInput) (output *api.DescribeVolumesOutput, err error) {
	err = c.client.Invoke(ctx, "describe-volumes", input, &output)
	return
}

func (c *Workspace) DescribeNetworks(ctx context.Context, input *api.DescribeNetworksInput) (output *api.DescribeNetworksOutput, err error) {
	err = c.client.Invoke(ctx, "describe-networks", input, &output)
	return
}

func (c *Workspace) ExportProcfile(ctx context.Context, input *api.ExportProcfileInput) (output *api.ExportProcfileOutput, err error) {
	err = c.client.Invoke(ctx, "export-procfile", input, &output)
	return
}

func (c *Workspace) ReadFile(ctx context.Context, input *api.ReadFileInput) (output *api.ReadFileOutput, err error) {
	err = c.client.Invoke(ctx, "read-file", input, &output)
	return
}

func (c *Workspace) WriteFile(ctx context.Context, input *api.WriteFileInput) (output *api.WriteFileOutput, err error) {
	err = c.client.Invoke(ctx, "write-file", input, &output)
	return
}

func (c *Workspace) BuildComponents(ctx context.Context, input *api.BuildComponentsInput) (output *api.BuildComponentsOutput, err error) {
	err = c.client.Invoke(ctx, "build-components", input, &output)
	return
}

func (c *Workspace) DescribeEnvironment(ctx context.Context, input *api.DescribeEnvironmentInput) (output *api.DescribeEnvironmentOutput, err error) {
	err = c.client.Invoke(ctx, "describe-environment", input, &output)
	return
}
