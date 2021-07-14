package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/deref/exo/core"
	"github.com/deref/exo/jsonutil"
)

func (provider *Provider) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	procDir := filepath.Join(provider.VarDir, input.ID)
	var state state
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	provider.refresh(&state)

	if state.Pid == 0 {
		var err error
		state, err = provider.start(ctx, procDir, input.Spec)
		if err != nil {
			return nil, err
		}
	}

	var output core.StartOutput
	output.State = jsonutil.MustMarshalString(state)
	return &output, nil
}

func (provider *Provider) start(ctx context.Context, procDir string, inputSpec string) (state, error) {
	var spec spec
	if err := jsonutil.UnmarshalString(inputSpec, &spec); err != nil {
		return state{}, fmt.Errorf("unmarshalling spec: %w", err)
	}

	// Use configured working directory or fallback to project directory.
	directory := spec.Directory
	if directory == "" {
		directory = provider.ProjectDir
	}

	// Resolve command path.
	command := spec.Command
	searchPaths, _ := os.LookupEnv("PATH")
	for _, searchPath := range strings.Split(searchPaths, ":") {
		candidate := filepath.Join(searchPath, command)
		info, _ := os.Stat(candidate)
		if info != nil {
			command = candidate
			break
		}
	}

	// Construct supervised command.
	fifofumPath := "./fifofum" // XXX Use exo home path.
	fifofumArgs := append(
		[]string{
			procDir,
			command,
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
		return state{}, fmt.Errorf("starting fifofum: %w", err)
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
	return state{Pid: pid}, err
}

func (provider *Provider) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	var state state
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}
	if state.Pid == 0 {
		return nil, nil
	}

	provider.stop(state.Pid)
	return &core.StopOutput{State: "null"}, nil
}

func (provider *Provider) stop(pid int) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	if err := proc.Kill(); err != nil {
		// TODO: Report the error somehow?
	}
}
