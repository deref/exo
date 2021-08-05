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

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/osutil"
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

	changed := make(chan string, 1)
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

		delay:
			for {
				timeout := time.After(50 * time.Millisecond)
				select {
				case <-changed:
					// Got additional change notifications - ignore.
				case <-timeout:
					// No additional notifications.
					break delay
				}
			}

			switch state {
			case runStateInitial:
				<-childStarted
				state = runStateRunning

			case runStateRunning:
				select {
				case name := <-changed:
					fmt.Printf("%q changed. Restarting process.\n", name)
					gracefullyShutdown(child, syscall.SIGTERM)
					state = runStateRestarting

				case sig := <-stopSignals:
					gracefullyShutdown(child, sig)
					state = runStateWaitingChildExit

				case <-childStopped:
					// Child stopped when we were not expecting a restart, so exit.
					done <- struct{}{}
				}

			case runStateRestarting:
				gracefulShutdownTimer := time.NewTimer(time.Second * time.Duration(10))
				select {
				case <-childStopped:
					gracefulShutdownTimer.Stop()

				case <-gracefulShutdownTimer.C:
					osutil.KillProcessGroup(child.Process.Pid)
				}
				startNewChild()
				state = runStateInitial

			case runStateWaitingChildExit:
				gracefulShutdownTimer := time.NewTimer(time.Second * time.Duration(10))
				select {
				case <-childStopped:
					gracefulShutdownTimer.Stop()

				case <-gracefulShutdownTimer.C:
					osutil.KillProcessGroup(child.Process.Pid)
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
			// Watch for file change events.
			case event := <-watcher.Events:
				changed <- event.Name

			// Exit on error.
			case err := <-watcher.Errors:
				if child != nil {
					gracefullyShutdown(child, syscall.SIGTERM)
				}
				cmdutil.Fatalf("watching files: %w", err)
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

				if err := watcher.Add(path); err != nil {
					return fmt.Errorf("watching %q: %w", path, err)
				}
				return nil
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
	if err := osutil.SignalProcessGroup(cmd.Process.Pid, sysSig); err != nil && err != os.ErrProcessDone {
		cmdutil.Fatalf("signaling child: %w", err)
	}
}
