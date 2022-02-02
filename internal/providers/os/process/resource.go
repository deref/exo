package process

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/deref/exo/internal/providers/sdk"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/which"
)

func (c *Controller) Identify(ctx context.Context, m *Model) (string, error) {
	if m.Pid == nil {
		return "", nil
	}
	return fmt.Sprintf("exo:/processes/%d", *m.Pid), nil
}

func (c *Controller) Create(ctx context.Context, m *Model) error {
	// Resolve program path.
	{
		whichQ := which.Query{
			Program: m.Program,
		}
		whichQ.WorkingDirectory = m.Directory
		whichQ.PathVariable = m.Environment["PATH"]
		var err error
		m.ProgramPath, err = whichQ.Run()
		if err != nil {
			return errutil.WithHTTPStatus(http.StatusBadRequest, err)
		}
	}

	// Run process.
	cmd := &exec.Cmd{
		Path: m.ProgramPath,
		Args: append([]string{m.Program}, m.Arguments...),
		Dir:  m.Directory,
		Env:  osutil.EnvMapToEnvv(m.Environment),
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

func (c *Controller) Shutdown(ctx context.Context, m *Model) error {
	if !c.Exists(ctx, m) {
		return sdk.ErrResourceGone
	}
	return osutil.ShutdownGroup(ctx, *m.Pid)
}

func (c *Controller) Delete(ctx context.Context, m *Model) error {
	if m.Pid == nil {
		return nil
	}
	return osutil.KillGroup(*m.Pid)
}
