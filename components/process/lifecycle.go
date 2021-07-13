package process

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

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

	// Forward environment.
	cmd.Env = []string{} // TODO: Get from spec.

	// Connect pipes.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start supervisor process.
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting fifofum: %w", err)
	}

	// Collect fifofum output.
	pidC := make(chan int, 1)
	errC := make(chan error, 2)
	go func() {
		pidStr, err := readLine(stdout)
		if err != nil {
			errC <- fmt.Errorf("reading fifofum stdout: %w", err)
			return
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			errC <- fmt.Errorf("parsing fifofum output: %w", err)
			return
		}
		pidC <- pid
	}()
	go func() {
		message, err := readLine(stderr)
		if err != nil {
			errC <- fmt.Errorf("reading fifofum stderr: %w", err)
			return
		}
		if len(message) > 0 {
			errC <- errors.New(message)
		}
	}()

	// Await fifofum result.
	var pid int
	select {
	case pid = <-pidC:
	case err = <-errC:
	case <-time.After(300 * time.Millisecond):
		err = errors.New("fifofum timeout")
	}
	if err != nil {
		return nil, err
	}

	var output api.InitializeOutput
	output.State = map[string]interface{}{
		"pid": pid,
	}
	return &output, nil
}

func readLine(r io.Reader) (string, error) {
	b := bufio.NewReader(r)
	line, isPrefix, err := b.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("line too long")
	}
	return string(line), nil
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
