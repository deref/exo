package process

import (
	"context"
	"encoding/json"
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
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/which"
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

	gracePeriod := 5
	if p.ShutdownGracePeriodSeconds != nil {
		gracePeriod = *p.ShutdownGracePeriodSeconds
	}

	childEnv, err := json.Marshal(p.Environment)
	if err != nil {
		return fmt.Errorf("encoding child environment")
	}

	// Construct supervised command.
	supervisePath := os.Args[0]
	superviseArgs := append(
		[]string{
			"supervise",
			"--",
			strconv.Itoa(int(p.SyslogPort)),
			p.ComponentID,
			p.WorkspaceRoot,
			strconv.Itoa(gracePeriod),
			string(childEnv),
			program,
		},
		p.Arguments...,
	)
	cmd := exec.Command(supervisePath, superviseArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Run in background.
	}

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
		p.Logger.Infof("interrupt failed: %w", err)
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
