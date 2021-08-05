package container

import (
	"context"
	"time"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
)

func (c *Container) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	if err := c.start(ctx); err != nil {
		return nil, err
	}
	return &core.StartOutput{}, nil
}

func (c *Container) start(ctx context.Context) error {
	err := c.Docker.ContainerStart(ctx, c.ContainerID, types.ContainerStartOptions{})
	if err != nil {
		c.Running = true
	}
	return err
}

func (c *Container) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	if err := c.stop(ctx); err != nil {
		return nil, err
	}
	return &core.StopOutput{}, nil
}

func (c *Container) stop(ctx context.Context) error {
	var timeout *time.Duration // Use container's default stop timeout.
	return c.Docker.ContainerStop(ctx, c.ContainerID, timeout)
}

func (c *Container) Restart(ctx context.Context, input *core.RestartInput) (*core.RestartOutput, error) {
	if err := c.restart(ctx); err != nil {
		return nil, err
	}
	return &core.RestartOutput{}, nil
}

func (c *Container) restart(ctx context.Context) error {
	var timeout *time.Duration // Use container's default stop timeout.
	return c.Docker.ContainerRestart(ctx, c.ContainerID, timeout)
}
