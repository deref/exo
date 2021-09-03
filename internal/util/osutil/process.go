package osutil

import (
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

func KillProcess(pid int) error {
	return SignalProcess(pid, syscall.SIGKILL)
}

func KillGroup(pgid int) error {
	return SignalGroup(pgid, syscall.SIGKILL)
}

func WaitProcess(pid int) (*os.ProcessState, error) {
	return FindProcess(pid).Wait()
}

func TerminateProcessWithTimeout(pid int, timeout time.Duration) error {
	_ = SignalProcess(pid, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		_, _ = WaitProcess(pid)
		close(done)
	}()

	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		return KillProcess(pid)
	case <-done:
		timer.Stop()
		return nil
	}
}

func TerminateGroupWithTimeout(pgid int, timeout time.Duration) error {
	return TerminateProcessWithTimeout(-pgid, timeout)
}
