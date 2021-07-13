package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deref/exo/api"
	"github.com/mitchellh/mapstructure"
)

type Lifecycle struct {
	ProjectDir string
	VarDir     string
}

type spec struct {
	Directory string
	Command   string
	Arguments []string
}

func (lc *Lifecycle) Initialize(ctx context.Context, input *api.InitializeInput) (*api.InitializeOutput, error) {
	var spec spec
	if err := mapstructure.Decode(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("destructuring spec: %w", err)
	}

	// Ensure top-level var directory.
	err := os.Mkdir(lc.VarDir, 0700)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("creating var directory: %w", err)
	}

	// Create var directory for the new process.
	procDir := filepath.Join(lc.VarDir, input.ID)
	if err := os.Mkdir(procDir, 0700); err != nil {
		return nil, fmt.Errorf("creating proc directory: %w", err)
	}

	// Use configured working directory or fallback to project directory.
	directory := spec.Directory
	if directory == "" {
		directory = lc.ProjectDir
	}

	// Forward environment.
	envv := []string{} // TODO: Get from spec.

	// Construct supervised command.
	fifofumPath := "./fifofum" // XXX Use exo home path.
	fifofumArgs := append(
		[]string{
			procDir,
			spec.Command,
		},
		spec.Arguments...,
	)
	cmd := exec.Command(fifofumPath, fifofumArgs...)

	// Start supervisor process.
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting: %w", err)
	}

	go func() {
		// XXX read stdin for pid.
	}()
	go func() {
		// XXX read stdout for errors.
	}()

	// XXX

	var output api.InitializeOutput
	output.State = map[string]interface{}{
		"pid": proc.Pid,
	}
	return &output, nil
}

func (lc *Lifecycle) Update(context.Context, *api.UpdateInput) (*api.UpdateOutput, error) {
	panic("TODO: update")
}

func (lc *Lifecycle) Refresh(context.Context, *api.RefreshInput) (*api.RefreshOutput, error) {
	panic("TODO: refresh")
}

func (lc *Lifecycle) Dispose(context.Context, *api.DisposeInput) (*api.DisposeOutput, error) {
	panic("TODO: dispose")
}
