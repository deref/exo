package container

import (
	"context"
	"fmt"
	"time"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
)

func (c *Container) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	if err := c.start(ctx); err != nil {
		return nil, fmt.Errorf("starting process container: %w", err)
	}
	return &core.StartOutput{}, nil
}

func (c *Container) start(ctx context.Context) error {
	err := c.Docker.ContainerStart(ctx, c.State.ContainerID, types.ContainerStartOptions{})
	if err != nil {
		c.State.Running = true
		return fmt.Errorf("starting container: %w", err)
	}
	return nil
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

	if err := c.Docker.ContainerStop(ctx, c.State.ContainerID, timeout); err != nil {
		//if strings.Contains(err.Error(), "No such container") {
		//return nil
		//}
		return fmt.Errorf("stopping container: %w", err)
	}
	return nil
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
	if err := c.Docker.ContainerRestart(ctx, c.State.ContainerID, timeout); err != nil {
		return fmt.Errorf("restarting container: %w", err)
	}
	return nil
}

func (c *Container) Signal(ctx context.Context, input *core.SignalInput) (*core.SignalOutput, error) {
	if err := c.Docker.ContainerKill(ctx, c.State.ContainerID, input.Signal); err != nil {
		return nil, fmt.Errorf("signaling container: %w", err)
	}
	return &core.SignalOutput{}, nil
}
