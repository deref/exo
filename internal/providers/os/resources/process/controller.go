package process

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/deref/exo/internal/providers/sdk"
	"github.com/deref/exo/internal/util/osutil"
)

type Controller struct{}

func (c *Controller) Identify(ctx context.Context, m *Model) (string, error) {
	if m.Pid == nil {
		return "", nil
	}
	return fmt.Sprintf("exo:/processes/%d", *m.Pid), nil
}

func (c *Controller) Create(ctx context.Context, m *Model) error {
	cmd := &exec.Cmd{
		Path: m.Program,
		Args: append([]string{m.Program}, m.Arguments...),
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	m.Pid = &cmd.Process.Pid
	return nil
}

func (c *Controller) Read(ctx context.Context, m *Model) error {
	if !c.Exists(ctx, m) {
		return sdk.ErrResourceGone
	}
	return nil
}

func (c *Controller) Exists(ctx context.Context, m *Model) bool {
	return m.Pid != nil && osutil.IsValidPid(*m.Pid)
}

func (c *Controller) Delete(ctx context.Context, m *Model) error {
	if m.Pid == nil {
		return nil
	}
	return osutil.KillGroup(*m.Pid)
}
