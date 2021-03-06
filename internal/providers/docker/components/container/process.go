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
	err := c.Docker.ContainerStart(ctx, c.State.ContainerID, types.ContainerStartOptions{})
	if err != nil {
		c.State.Running = true
	}
	return err
}

func (c *Container) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	if c.State.ContainerID == "" {
		return &core.StopOutput{}, nil
	}
	if err := c.stop(ctx, input.TimeoutSeconds); err != nil {
		return nil, err
	}
	return &core.StopOutput{}, nil
}

func (c *Container) stop(ctx context.Context, timeoutSeconds *uint) error {
	var timeout *time.Duration // Use container's default stop timeout.
	if timeoutSeconds != nil {
		duration := time.Second * time.Duration(*timeoutSeconds)
		timeout = &duration
	}

	return c.Docker.ContainerStop(ctx, c.State.ContainerID, timeout)
}

func (c *Container) Restart(ctx context.Context, input *core.RestartInput) (*core.RestartOutput, error) {
	if err := c.restart(ctx, input.TimeoutSeconds); err != nil {
		return nil, err
	}
	return &core.RestartOutput{}, nil
}

func (c *Container) restart(ctx context.Context, timeoutSeconds *uint) error {
	var timeout *time.Duration // Use container's default stop timeout.
	if timeoutSeconds != nil {
		duration := time.Second * time.Duration(*timeoutSeconds)
		timeout = &duration
	}
	return c.Docker.ContainerRestart(ctx, c.State.ContainerID, timeout)
}

func (c *Container) Signal(ctx context.Context, input *core.SignalInput) (*core.SignalOutput, error) {
	if err := c.Docker.ContainerKill(ctx, c.State.ContainerID, input.Signal); err != nil {
		return nil, err
	}
	return &core.SignalOutput{}, nil
}
