package process

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	core "github.com/deref/exo/core/api"
	"github.com/deref/exo/util/errutil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/which"
)

func (provider *Provider) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	procDir := filepath.Join(provider.VarDir, input.ID)
	var state State
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

func (provider *Provider) start(ctx context.Context, procDir string, inputSpec string) (State, error) {
	var spec Spec
	if err := jsonutil.UnmarshalString(inputSpec, &spec); err != nil {
		return State{}, fmt.Errorf("unmarshalling spec: %w", err)
	}

	// Use configured working directory or fallback to workspace directory.
	whichQ := which.Query{
		Program: spec.Program,
	}
	whichQ.WorkingDirectory = spec.Directory
	if whichQ.WorkingDirectory == "" {
		whichQ.WorkingDirectory = provider.WorkspaceDir
	}
	whichQ.PathVariable = spec.Environment["PATH"]
	if whichQ.PathVariable == "" {
		// TODO: Daemon path from config.
		whichQ.PathVariable, _ = os.LookupEnv("PATH")
	}
	program, err := whichQ.Run()
	if err != nil {
		return State{}, errutil.WithHTTPStatus(http.StatusBadRequest, err)
	}

	// Construct supervised command.
	fifofumPath := os.Args[0]
	fifofumArgs := append(
		[]string{
			"fifofum",
			procDir,
			program,
		},
		spec.Arguments...,
	)
	cmd := exec.Command(fifofumPath, fifofumArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Run in background.
	}

	// Forward environment.
	envv := os.Environ()
	envMap := make(map[string]string, len(envv)+len(spec.Environment))
	addEnv := func(key, val string) {
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		envMap[key] = val
	}
	for _, kvp := range envv {
		parts := strings.SplitN(kvp, "=", 2)
		if len(parts) < 2 {
			parts = append(parts, "")
		}
		addEnv(parts[0], parts[1])
	}
	for key, val := range spec.Environment {
		addEnv(key, val)
	}
	for key, val := range envMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}
	sort.Strings(cmd.Env)

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
		return State{}, fmt.Errorf("starting fifofum: %w", err)
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
			// TODO: Do not treat as a bad request. Record the error somewhere,
			// mark the component as being in an error state.
			errC <- errutil.NewHTTPError(http.StatusBadRequest, message)
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
	return State{Pid: pid}, err
}

func (provider *Provider) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	var state State
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
