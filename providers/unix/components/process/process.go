package process

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	core "github.com/deref/exo/core/api"
	"github.com/deref/exo/util/errutil"
	"github.com/deref/exo/util/which"
)

func (p *Process) Start(ctx context.Context, input *core.StartInput) (*core.StartOutput, error) {
	p.refresh()
	if p.Pid == 0 {
		if err := p.start(ctx); err != nil {
			return nil, err
		}
	}
	return &core.StartOutput{}, nil
}

func (p *Process) start(ctx context.Context) error {
	whichQ := which.Query{
		Program: p.Program,
	}
	whichQ.WorkingDirectory = p.Directory
	if whichQ.WorkingDirectory == "" {
		whichQ.WorkingDirectory = p.WorkspaceDir
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

	gracePeriod := 5
	if p.ShutdownGracePeriodSeconds != nil {
		gracePeriod = *p.ShutdownGracePeriodSeconds
	}

	// Construct supervised command.
	supervisePath := os.Args[0]
	superviseArgs := append(
		[]string{
			"supervise",
			"--",
			strconv.Itoa(p.SyslogPort),
			p.ComponentID,
			p.WorkspaceDir,
			strconv.Itoa(gracePeriod),
			program,
		},
		p.Arguments...,
	)
	cmd := exec.Command(supervisePath, superviseArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Run in background.
	}

	// Forward environment.
	envv := os.Environ()
	envMap := make(map[string]string, len(envv)+len(p.Environment))
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
	for key, val := range p.Environment {
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
		return fmt.Errorf("starting supervise: %w", err)
	}
	p.State.SupervisorPid = cmd.Process.Pid
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
		if len(message) > 0 {
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
	p.stop()
	return &core.StopOutput{}, nil
}

func (p *Process) stop() {
	if p.SupervisorPid == 0 {
		return
	}
	proc, err := os.FindProcess(p.SupervisorPid)
	if err != nil {
		panic(err)
	}
	p.Pid = 0
	if err := proc.Signal(os.Interrupt); err != nil {
		// TODO: Report the error somehow?
	}
}

func (p *Process) Restart(ctx context.Context, input *core.RestartInput) (*core.RestartOutput, error) {
	p.stop()
	err := p.start(ctx)
	if err != nil {
		return nil, err
	}
	return &core.RestartOutput{}, nil
}
