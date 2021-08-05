package osutil

import (
	"os"
	"syscall"
	"time"
)

func IsValidPid(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func SignalProcessGroup(pgid int, sig syscall.Signal) error {
	return syscall.Kill(-pgid, sig)
}

func KillProcessGroup(pgid int) error {
	return SignalProcessGroup(pgid, syscall.SIGKILL)
}

func FindProcess(pid int) *os.Process {
	proc, err := os.FindProcess(pid)
	if err != nil {
		// Should be unreachable on all supported operating systems.
		panic(err)
	}
	return proc
}

func TerminateProcessGroupWithTimeout(pgid int, timeout time.Duration) error {
	proc := FindProcess(pgid)
	_ = proc.Signal(syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		_, _ = proc.Wait()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		return KillProcessGroup(pgid)
	case <-done:
		timer.Stop()
		return nil
	}
}
