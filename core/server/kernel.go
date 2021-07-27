package server

import (
	"context"

	"github.com/deref/exo/core/api"
	state "github.com/deref/exo/core/state/api"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/telemetry"
)

type Kernel struct {
	VarDir string
	Store  state.Store
}

func (kern *Kernel) CreateWorkspace(ctx context.Context, input *api.CreateWorkspaceInput) (*api.CreateWorkspaceOutput, error) {
	id := gensym.RandomBase32()
	_, err := kern.Store.AddWorkspace(ctx, &state.AddWorkspaceInput{
		ID:   id,
		Root: input.Root,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateWorkspaceOutput{
		ID: id,
	}, nil
}

func (kern *Kernel) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (*api.DescribeWorkspacesOutput, error) {
	output, err := kern.Store.DescribeWorkspaces(ctx, &state.DescribeWorkspacesInput{})
	if err != nil {
		return nil, err
	}
	workspaces := make([]api.WorkspaceDescription, len(output.Workspaces))
	for i, workspace := range output.Workspaces {
		workspaces[i] = api.WorkspaceDescription{
			ID:   workspace.ID,
			Root: workspace.Root,
		}
	}
	return &api.DescribeWorkspacesOutput{
		Workspaces: workspaces,
	}, nil
}

func (kern *Kernel) FindWorkspace(ctx context.Context, input *api.FindWorkspaceInput) (*api.FindWorkspaceOutput, error) {
	output, err := kern.Store.FindWorkspace(ctx, &state.FindWorkspaceInput{
		Path: input.Path,
	})
	if err != nil {
		return nil, err
	}
	return &api.FindWorkspaceOutput{
		ID: output.ID,
	}, nil
}

func (kern *Kernel) GetVersion(ctx context.Context, input *api.GetVersionInput) (*api.GetVersionOutput, error) {
	installed := telemetry.CurrentVersion()
	current := true
	var latest *string
	if telemetry.CanSelfUpgrade() {
		latestVersion, err := telemetry.LatestVersion()
		if err != nil {
			return nil, err
		}
		latest = &latestVersion
		current = installed >= latestVersion
	}

	return &api.GetVersionOutput{
		Installed: installed,
		Latest:    latest,
		Current:   current,
	}, nil
}

func (kern *Kernel) Panic(ctx context.Context, input *api.PanicInput) (*api.PanicOutput, error) {
	message := input.Message
	if input.Message == "" {
		message = "test error"
	}
	panic(message)
}
