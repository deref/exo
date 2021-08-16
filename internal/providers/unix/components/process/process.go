package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/supervise"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/which"
)

func (p *Process) zeroPids() bool {
	return p.Pgid == 0 && p.SupervisorPid == 0 && p.Pid == 0
}

func (p *Process) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	p.refresh()
	if p.zeroPids() {
		if err := p.start(ctx); err != nil {
			return nil, err
		}
	}
	return &core.StartOutput{}, nil
}

func (p *Process) start(ctx context.Context) error {
	p.State.clear()

	whichQ := which.Query{
		Program: p.Program,
	}
	whichQ.WorkingDirectory = p.Directory
	if whichQ.WorkingDirectory == "" {
		whichQ.WorkingDirectory = p.WorkspaceRoot
	}
	whichQ.PathVariable = p.Environment["PATH"]
	if whichQ.PathVariable == "" {
		// TODO: Daemon path from config.
		whichQ.PathVariable, _ = os.LookupEnv("PATH")
	}
	program, err := whichQ.Run()
	if err != nil {
		return errutil.WithHTTPStatus(http.StatusBadRequest, err)
	}

	// Construct supervised command.
	supervisePath := os.Args[0]
	cmd := exec.Command(supervisePath, "supervise")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Run in background.
	}

	// Pipe JSON config to supervise on stdin.
	configJSON := supervise.MustEncodeConfig(&supervise.Config{
		ComponentID:      p.ComponentID,
		WorkingDirectory: p.WorkspaceRoot,
		SyslogPort:       p.SyslogPort,
		Environment:      p.Environment,
		Program:          program,
		Arguments:        p.Arguments,
	})
	cmd.Stdin = bytes.NewBuffer(configJSON)

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
		return fmt.Errorf("starting supervise: %w", err)
	}
	p.State.SupervisorPid = cmd.Process.Pid
	p.State.Pgid, _ = syscall.Getpgid(p.State.SupervisorPid)
	p.State.Pid = 0 // Overriden below.

	envMap := make(map[string]string)
	for _, assign := range os.Environ() {
		parts := strings.SplitN(assign, "=", 2)
		key := parts[0]
		val := parts[1]
		envMap[key] = val
	}
	for key, val := range p.Environment {
		envMap[key] = val
	}
	p.State.FullEnvironment = envMap

	// Collect supervise output.
	pidC := make(chan int, 1)
	errC := make(chan error, 2)
	go func() {
		pidStr, err := readLine(stdout)
		if err != nil {
			errC <- fmt.Errorf("reading supervise stdout: %w", err)
			return
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			errC <- fmt.Errorf("parsing supervise output: %w", err)
			return
		}
		pidC <- pid
	}()
	go func() {
		message, err := readLine(stderr)
		if err != nil {
			errC <- fmt.Errorf("reading supervise stderr: %w", err)
			return
		}
		// SEE NOTE: [SUPERVISE_STDERR].
		if len(message) > 0 && message != "started ok" {
			// TODO: Do not treat as a bad request. Record the error somewhere,
			// mark the component as being in an error state.
			errC <- errutil.NewHTTPError(http.StatusBadRequest, message)
		}
	}()

	// Await supervise result.
	select {
	case p.Pid = <-pidC:
	case err = <-errC:
	case <-time.After(300 * time.Millisecond):
		err = errors.New("supervise startup timeout")
	}
	return err
}

func (p *Process) Stop(ctx context.Context, input *core.StopInput) (*core.StopOutput, error) {
	p.stop(input.StopNow)
	return &core.StopOutput{}, nil
}

const DefaultShutdownGracePeriod = 5 * time.Second

func (p *Process) stop(stopNow bool) {
	if p.zeroPids() {
		return
	}

	timeout := DefaultShutdownGracePeriod
	if p.ShutdownGracePeriodSeconds != nil {
		timeout = time.Duration(*p.ShutdownGracePeriodSeconds) * time.Second
	}
	if stopNow {
		timeout = time.Duration(0)
	}

	if err := osutil.TerminateGroupWithTimeout(p.Pgid, timeout); err != nil {
		p.Logger.Infof("terminating process: %w", err)
	}

	p.State.clear()
}

func (p *Process) Restart(ctx context.Context, input *core.RestartInput) (*core.RestartOutput, error) {
	p.stop(false)
	err := p.start(ctx)
	if err != nil {
		return nil, err
	}
	return &core.RestartOutput{}, nil
}
