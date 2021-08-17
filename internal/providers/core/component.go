package core

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/logging"
)

type Component interface {
	GetComponentID() string

	// TODO: Rethink these.
	IsDeleted() bool
	MarkDeleted()
}

type ComponentBase struct {
	ComponentID          string
	WorkspaceRoot        string
	WorkspaceEnvironment map[string]string
	Logger               logging.Logger
	isDeleted            bool
	ExoLabels            map[string]string
}

func (c ComponentBase) GetComponentID() string {
	return c.ComponentID
}

func (c *ComponentBase) IsDeleted() bool {
	return c.isDeleted
}

func (c *ComponentBase) MarkDeleted() {
	c.isDeleted = true
}

func (c *ComponentBase) Build(ctx context.Context, input *api.BuildInput) (*api.BuildOutput, error) {
	// Default no-op implemention.
	return &api.BuildOutput{}, nil
}
