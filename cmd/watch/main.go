package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/deref/exo/util/cmdutil"
	"github.com/fsnotify/fsnotify"
)

type runState int

const (
	runStateInitial runState = iota
	runStateRunning
	runStateRestarting
	runStateWaitingChildExit
)

func (s runState) String() string {
	switch s {
	case runStateInitial:
		return "Initial"
	case runStateRunning:
		return "Running"
	case runStateRestarting:
		return "Restarting"
	case runStateWaitingChildExit:
		return "WaitingChildExit"
	}
	return "<invalid state>"
}

func main() {
	cmd, err := cmdutil.ParseArgs(os.Args)
	if err != nil {
		cmdutil.Fatalf("parsing arguments: %w", err)
	}

	if len(cmd.Args) == 0 {
		cmdutil.Fatalf("expected command to execute")
	}

	reload := make(chan struct{}, 1)
	childStarted := make(chan struct{}, 1)
	childStopped := make(chan struct{}, 1)
	stopSignals := make(chan os.Signal)
	done := make(chan struct{})

	var child *exec.Cmd
	startNewChild := func() {
		var err error
		if child, err = startCommand(cmd.Args, childStarted, childStopped); err != nil {
			cmdutil.Fatalf("starting child process: %w", err)
		}
	}

	signal.Notify(stopSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Main loop
	go func() {
		state := runStateInitial
		for {
			select {
			case sig := <-stopSignals:
				gracefullyShutdown(child, sig)
				state = runStateWaitingChildExit

			default:
				// Did not receive a stop signal since last iteration - continue.
			}

			select {
			case _ = <-reload:
				// Got additional reload notifications - ignore.
			default:
				// No additional notifications.
			}

			switch state {
			case runStateInitial:
				<-childStarted
				state = runStateRunning

			case runStateRunning:
				select {
				case _ = <-reload:
					fmt.Println("Files changed. Restarting process.")
					gracefullyShutdown(child, syscall.SIGTERM)
					state = runStateRestarting

				case sig := <-stopSignals:
					gracefullyShutdown(child, sig)
					state = runStateWaitingChildExit

				case _ = <-childStopped:
					// Child stopped when we were not expecting a restart, so exit.
					done <- struct{}{}
				}

			case runStateRestarting:
				gracefulShutdownTimer := time.NewTimer(time.Second * time.Duration(10))
				select {
				case _ = <-childStopped:
					gracefulShutdownTimer.Stop()

				case _ = <-gracefulShutdownTimer.C:
					child.Process.Kill()
					// TODO: try to clean up other processes in the process group.
				}
				startNewChild()
				state = runStateInitial

			case runStateWaitingChildExit:
				gracefulShutdownTimer := time.NewTimer(time.Second * time.Duration(10))
				select {
				case _ = <-childStopped:
					gracefulShutdownTimer.Stop()

				case _ = <-gracefulShutdownTimer.C:
					child.Process.Kill()
					// TODO: try to clean up other processes in the process group.
				}
				done <- struct{}{}
			}
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		cmdutil.Fatalf("creating watcher: %w", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			// watch for events
			case _ = <-watcher.Events:
				reload <- struct{}{}

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	dir, ok := cmd.Flags["dir"]
	if !ok {
		dir = "."
	}

	if _, ok := cmd.Flags["r"]; ok {
		var ignores []string
		if ignoreStr, ok := cmd.Flags["ignore"]; ok {
			ignores = strings.Split(ignoreStr, ",")
		}

		if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if info.Mode().IsDir() {
				// TODO: implement something a little more robust than substring match for ignore.
				for _, ignore := range ignores {
					if strings.Contains(path, ignore) {
						return nil
					}
				}

				return watcher.Add(path)
			}

			return nil
		}); err != nil {
			cmdutil.Fatalf("walk directory for watches: %w", err)
		}
	} else {
		err = watcher.Add(dir)
	}

	startNewChild()
	<-done
}

func startCommand(invocation []string, started, stopped chan struct{}) (*exec.Cmd, error) {
	program := invocation[0]
	args := invocation[1:]

	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Assign child a new process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	started <- struct{}{}

	go func() {
		_ = cmd.Wait()
		// TODO: communicate exit code.
		stopped <- struct{}{}
	}()

	return cmd, nil
}

func gracefullyShutdown(cmd *exec.Cmd, sig os.Signal) {
	sysSig := sig.(syscall.Signal)
	if err := syscall.Kill(-cmd.Process.Pid, sysSig); err != nil && err != os.ErrProcessDone {
		cmdutil.Fatalf("signaling child: %w", err)
	}
}
