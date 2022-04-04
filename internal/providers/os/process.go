package os

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/which"
	"github.com/deref/exo/sdk"
)

type ProcessController struct {
	sdk.ResourceController[ProcessModel]
}

func NewProcessController(svc api.Service) *sdk.ResourceComponentController {
	return sdk.NewResourceComponentController[ProcessModel](svc, &ProcessController{})
}

type ProcessModel struct {
	ProcessSpec
	ProcessState
}

type ProcessSpec struct {
	Program     string            `json:"program,omitempty"`
	Arguments   []string          `json:"arguments,omitempty"`
	Directory   string            `json:"directory,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

type ProcessState struct {
	ProgramPath string `json:"programPath,omitempty"`
	Pid         *int   `json:"pid,omitempty"`
}

func (ctrl *ProcessController) IdentifyResource(ctx context.Context, cfg *sdk.ResourceConfig, m *ProcessModel) (string, error) {
	if m.Pid == nil {
		return "", nil
	}
	return fmt.Sprintf("exo:/processes/%d", *m.Pid), nil
}

func (ctrl *ProcessController) Create(ctx context.Context, cfg *sdk.ResourceConfig, m *ProcessModel) error {
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
		SysProcAttr: &syscall.SysProcAttr{
			Setsid: true,
		},
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	m.Pid = &cmd.Process.Pid

	return nil
}

func (ctrl *ProcessController) ReadResource(ctx context.Context, cfg *sdk.ResourceConfig, m *ProcessModel) error {
	if !ctrl.exists(m) {
		return sdk.ErrResourceGone
	}
	return nil
}

func (ctrl *ProcessController) Shutdown(ctx context.Context, cfg *sdk.ResourceConfig, m *ProcessModel) error {
	if !ctrl.exists(m) {
		return sdk.ErrResourceGone
	}
	return osutil.ShutdownGroup(ctx, *m.Pid)
}

func (ctrl *ProcessController) exists(m *ProcessModel) bool {
	return m.Pid != nil && osutil.IsValidPid(*m.Pid)
}

func (ctrl *ProcessController) Delete(ctx context.Context, cfg *sdk.ResourceConfig, m *ProcessModel) error {
	if m.Pid == nil {
		return nil
	}
	return osutil.KillGroup(*m.Pid)
}
