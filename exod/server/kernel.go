package server

import (
	"context"

	"github.com/deref/exo/exod/api"
	state "github.com/deref/exo/exod/state/api"
	"github.com/deref/exo/gensym"
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

func (kern *Kernel) Panic(ctx context.Context, input *api.PanicInput) (*api.PanicOutput, error) {
	message := input.Message
	if input.Message == "" {
		message = "test error"
	}
	panic(message)
}
