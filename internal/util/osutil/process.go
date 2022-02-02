package osutil

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"
)

func FindProcess(pid int) *os.Process {
	proc, err := os.FindProcess(pid)
	if err != nil {
		// Should be unreachable on all supported operating systems.
		panic(err)
	}
	return proc
}

func IsValidPid(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc := FindProcess(pid)
	err := proc.Signal(syscall.Signal(0))
	return err == nil
}

func SignalProcess(pid int, sig os.Signal) error {
	proc := FindProcess(pid)
	return proc.Signal(sig)
}

func SignalGroup(pgid int, sig os.Signal) error {
	return SignalProcess(-pgid, sig)
}

func TerminateProcess(pid int) error {
	return SignalProcess(pid, syscall.SIGTERM)
}

func TerminateGroup(pgid int) error {
	return SignalGroup(pgid, syscall.SIGTERM)
}

func KillProcess(pid int) error {
	return SignalProcess(pid, syscall.SIGKILL)
}

func KillGroup(pgid int) error {
	return SignalGroup(pgid, syscall.SIGKILL)
}

func WaitProcess(pid int) (*os.ProcessState, error) {
	return FindProcess(pid).Wait()
}

// Terminates pid, then if context is cancelled, kills pid.
func ShutdownProcess(ctx context.Context, pid int) error {
	if err := TerminateGroup(pid); err != nil {
		return fmt.Errorf("terminating: %w", err)
	}
	done := make(chan struct{})
	go func() {
		_, _ = WaitProcess(pid)
		close(done)
	}()
	select {
	case <-ctx.Done():
		if err := KillGroup(pid); err != nil {
			return fmt.Errorf("killing: %w", err)
		}
	case <-done:
		// no-op.
	}
	return nil
}

func ShutdownGroup(ctx context.Context, gpid int) error {
	return ShutdownProcess(ctx, -gpid)
}

func TerminateProcessWithTimeout(pid int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ShutdownProcess(ctx, pid)
}

func TerminateGroupWithTimeout(pgid int, timeout time.Duration) error {
	return TerminateProcessWithTimeout(-pgid, timeout)
}
