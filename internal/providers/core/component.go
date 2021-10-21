package core

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/logging"
)

type Component interface {
	GetComponentID() string
	GetComponentName() string
}

type ComponentBase struct {
	ComponentID          string
	ComponentName        string
	ComponentState       string
	WorkspaceID          string
	WorkspaceRoot        string
	WorkspaceEnvironment map[string]string
	StackName            string
	Logger               logging.Logger
	isDeleted            bool
}

func (c ComponentBase) GetComponentID() string {
	return c.ComponentID
}

func (c ComponentBase) GetComponentName() string {
	return c.ComponentName
}

func (c *ComponentBase) Build(ctx context.Context, input *api.BuildInput) (*api.BuildOutput, error) {
	// Default no-op implemention.
	return &api.BuildOutput{}, nil
}
