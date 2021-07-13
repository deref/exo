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
	"strings"
	"time"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/jsonutil"
	"github.com/deref/exo/kernel/api"
	logcol "github.com/deref/exo/logcol/api"
)

type Lifecycle struct {
	ProjectDir string
	VarDir     string
}

type spec struct {
	Directory string   `json:"directory"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

type state struct {
	Pid int `json:"pid"`
}

func (lc *Lifecycle) Initialize(ctx context.Context, input *api.InitializeInput) (*api.InitializeOutput, error) {
	var spec spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
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

	// Register logs.
	// TODO: Don't do this synchronously here. Use some kind of component hierarchy mechanism.
	collector := log.CurrentLogCollector(ctx)
	for _, role := range []string{"out", "err"} {
		_, err := collector.AddLog(ctx, &logcol.AddLogInput{
			Name:   fmt.Sprintf("%s:%s", input.ID, role),
			Source: filepath.Join(procDir, role),
		})
		if err != nil {
			return nil, fmt.Errorf("adding std%s log: %w", role, err)
		}
	}

	var output api.InitializeOutput
	output.State = jsonutil.MustMarshalString(state{
		Pid: pid,
	})
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

func (lc *Lifecycle) Dispose(ctx context.Context, input *api.DisposeInput) (*api.DisposeOutput, error) {
	var state state
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	proc, err := os.FindProcess(state.Pid)
	if err != nil {
		panic(err)
	}
	if err := proc.Kill(); err != nil {
		// TODO: Report the error somehow?
	}

	// Deregister log streams.
	// TODO: Don't do this synchronously here. Use some kind of component hierarchy mechanism.
	collector := log.CurrentLogCollector(ctx)
	for _, role := range []string{"out", "err"} {
		_, err := collector.RemoveLog(ctx, &logcol.RemoveLogInput{
			Name: fmt.Sprintf("%s:%s", input.ID, role),
		})
		if err != nil {
			return nil, fmt.Errorf("removing std%s log: %w", role, err)
		}
	}

	return &api.DisposeOutput{State: input.State}, nil
}
